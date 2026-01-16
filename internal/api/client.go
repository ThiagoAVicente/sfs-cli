package api

import (
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/vcnt/sfs-cli/internal/config"
)

// Client wraps the API client
type Client struct {
	client *resty.Client
	config *config.Config
}

// SearchResult represents a search result from the API
type SearchResult struct {
	Score   float64 `json:"score"`
	Payload struct {
		FilePath   string `json:"file_path"`
		Text       string `json:"text"`
		Start      int    `json:"start"`
		End        int    `json:"end"`
		ChunkIndex int    `json:"chunk_index"`
	} `json:"payload"`
}

// SearchResponse represents the search API response
type SearchResponse struct {
	Results []SearchResult `json:"results"`
}

// UploadResponse represents the upload API response
type UploadResponse struct {
	JobID string `json:"job_id"`
}

// JobStatusResponse represents the job status API response
type JobStatusResponse struct {
	JobID  string `json:"job_id"`
	Status string `json:"status"`
}

// ListFilesResponse represents the list files API response
type ListFilesResponse struct {
	Files []string `json:"files"`
	Count int      `json:"count"`
}

// DeleteResponse represents the delete API response
type DeleteResponse struct {
	JobID string `json:"job_id"`
}

// NewClient creates a new API client
func NewClient() (*Client, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("API key not configured. Run: sfs config set api_key <your-key>")
	}

	// Only skip cert validation for localhost
	isLocalhost := strings.Contains(cfg.APIURL, "://localhost") || strings.Contains(cfg.APIURL, "://127.0.0.1")

	client := resty.New().
		SetBaseURL(cfg.APIURL).
		SetHeader("X-API-Key", cfg.APIKey).
		SetHeader("Content-Type", "application/json").
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: isLocalhost})

	return &Client{
		client: client,
		config: cfg,
	}, nil
}

// UploadFile uploads a file to the API
func (c *Client) UploadFile(filePath string, update bool) (*UploadResponse, error) {
	// Convert to absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Convert file path to file name (replace / with _)
	// Remove leading / or \ and replace all path separators with _
	fileName := absPath
	if len(fileName) > 0 && (fileName[0] == '/' || fileName[0] == '\\') {
		fileName = fileName[1:]
	}
	fileName = replacePathSeparators(fileName)

	resp, err := c.client.R().
		SetFile("file", absPath).
		SetFormData(map[string]string{
			"update": fmt.Sprintf("%t", update),
		}).
		SetResult(&UploadResponse{}).
		Post("/index")

	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("upload failed: %s", resp.String())
	}

	result := resp.Result().(*UploadResponse)
	fmt.Printf("File uploaded: %s -> %s\n", absPath, fileName)
	fmt.Printf("Job ID: %s\n", result.JobID)

	return result, nil
}

// Search performs a semantic search
func (c *Client) Search(query string, limit int, scoreThreshold float64) (*SearchResponse, error) {
	body := map[string]interface{}{
		"query":           query,
		"limit":           limit,
		"score_threshold": scoreThreshold,
	}

	resp, err := c.client.R().
		SetBody(body).
		SetResult(&SearchResponse{}).
		Post("/search")

	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("search failed: %s", resp.String())
	}

	return resp.Result().(*SearchResponse), nil
}

// ListFiles lists all files, optionally filtered by prefix
func (c *Client) ListFiles(prefix string) (*ListFilesResponse, error) {
	req := c.client.R().SetResult(&ListFilesResponse{})

	if prefix != "" {
		req.SetQueryParam("prefix", prefix)
	}

	resp, err := req.Get("/files/")

	if err != nil {
		return nil, fmt.Errorf("list files failed: %w", err)
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("list files failed: %s", resp.String())
	}

	return resp.Result().(*ListFilesResponse), nil
}

// DeleteFile deletes a file
func (c *Client) DeleteFile(fileName string) (*DeleteResponse, error) {
	resp, err := c.client.R().
		SetResult(&DeleteResponse{}).
		Delete("/index/" + fileName)

	if err != nil {
		return nil, fmt.Errorf("delete failed: %w", err)
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("delete failed: %s", resp.String())
	}

	return resp.Result().(*DeleteResponse), nil
}

// DownloadFile downloads a file
func (c *Client) DownloadFile(fileName, destPath string) error {
	resp, err := c.client.R().
		SetDoNotParseResponse(true).
		Get("/files/" + fileName)

	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.RawBody().Close()

	if !resp.IsSuccess() {
		return fmt.Errorf("download failed: status %d", resp.StatusCode())
	}

	// Create destination file
	outFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	// Copy response body to file
	_, err = io.Copy(outFile, resp.RawBody())
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// GetJobStatus gets the status of an indexing job
func (c *Client) GetJobStatus(jobID string) (*JobStatusResponse, error) {
	resp, err := c.client.R().
		SetResult(&JobStatusResponse{}).
		Get("/index/status/" + jobID)

	if err != nil {
		return nil, fmt.Errorf("status check failed: %w", err)
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("status check failed: %s", resp.String())
	}

	return resp.Result().(*JobStatusResponse), nil
}

// replacePathSeparators replaces all path separators (/ and \) with underscores
func replacePathSeparators(path string) string {
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, "\\", "_")
	return path
}
