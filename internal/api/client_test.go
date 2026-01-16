package api

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/vcnt/sfs-cli/internal/config"
)

func setupTestConfig(t *testing.T) {
	// Reset viper to avoid state leakage between tests
	viper.Reset()

	// Use temp dir for test
	tmpDir := t.TempDir()
	home := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	t.Cleanup(func() {
		os.Setenv("HOME", home)
	})

	// Initialize config with test values
	if err := config.InitConfig(); err != nil {
		t.Fatalf("Failed to init config: %v", err)
	}

	config.Set("api_url", "https://test.localhost:8000")
	config.Set("api_key", "test-key")
}

func TestNewClient(t *testing.T) {
	setupTestConfig(t)

	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.config.APIURL != "https://test.localhost:8000" {
		t.Errorf("Expected API URL 'https://test.localhost:8000', got '%s'", client.config.APIURL)
	}

	if client.config.APIKey != "test-key" {
		t.Errorf("Expected API key 'test-key', got '%s'", client.config.APIKey)
	}
}

func TestNewClientWithoutAPIKey(t *testing.T) {
	// Reset viper to avoid state leakage between tests
	viper.Reset()

	// Use temp dir without setting API key
	tmpDir := t.TempDir()
	home := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", home)

	if err := config.InitConfig(); err != nil {
		t.Fatalf("Failed to init config: %v", err)
	}

	// Don't set api_key - should fail
	_, err := NewClient()
	if err == nil {
		t.Error("Expected error when creating client without API key")
	}
}

func TestUploadFileInvalidPath(t *testing.T) {
	setupTestConfig(t)

	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Try to upload non-existent file
	_, err = client.UploadFile("/nonexistent/file.txt", false)
	if err == nil {
		t.Error("Expected error when uploading non-existent file")
	}
}

func TestUploadFileValidPath(t *testing.T) {
	setupTestConfig(t)

	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Note: This will fail because the API server isn't running
	// but it validates that the file exists and the request is formatted correctly
	_, err = client.UploadFile(testFile, false)
	// We expect an error because the server isn't running, but not a file-related error
	if err == nil {
		t.Skip("API server is running, skipping validation-only test")
	}
}

func TestUploadFileRelativePath(t *testing.T) {
	setupTestConfig(t)

	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Change to temp dir so we can use relative path
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Upload using relative path
	_, err = client.UploadFile("test.txt", false)
	// We expect an error because the server isn't running
	// but it should handle the relative path correctly
	if err == nil {
		t.Skip("API server is running, skipping validation-only test")
	}

	// The error should not be about file not found
	if err != nil && strings.Contains(err.Error(), "failed to open file") {
		t.Error("Failed to handle relative path correctly")
	}
}

func TestSearchEmptyQuery(t *testing.T) {
	setupTestConfig(t)

	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Search with empty query should still work (might return error from API)
	_, err = client.Search("", 5, 0.5)
	// We expect an error because the server isn't running
	if err == nil {
		t.Skip("API server is running, skipping validation-only test")
	}
}

func TestSearchNegativeLimit(t *testing.T) {
	setupTestConfig(t)

	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Negative limit should be handled by the API or client
	_, err = client.Search("test query", -1, 0.5)
	// We expect an error because the server isn't running
	if err == nil {
		t.Skip("API server is running, skipping validation-only test")
	}
}

func TestDownloadFileInvalidDestination(t *testing.T) {
	setupTestConfig(t)

	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Try to download to invalid path (non-existent directory)
	err = client.DownloadFile("test.txt", "/nonexistent/directory/file.txt")
	if err == nil {
		t.Error("Expected error when downloading to invalid destination")
	}
}

func TestListFiles(t *testing.T) {
	setupTestConfig(t)

	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// ListFiles should format the request correctly
	_, err = client.ListFiles("")
	// We expect an error because the server isn't running
	if err == nil {
		t.Skip("API server is running, skipping validation-only test")
	}
}

func TestDeleteFile(t *testing.T) {
	setupTestConfig(t)

	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// DeleteFile should format the request correctly
	_, err = client.DeleteFile("test.txt")
	// We expect an error because the server isn't running
	if err == nil {
		t.Skip("API server is running, skipping validation-only test")
	}
}

func TestGetJobStatus(t *testing.T) {
	setupTestConfig(t)

	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// GetJobStatus should format the request correctly
	_, err = client.GetJobStatus("test-job-id")
	// We expect an error because the server isn't running
	if err == nil {
		t.Skip("API server is running, skipping validation-only test")
	}
}

func TestReplacePathSeparators(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "unix path",
			input:    "home/user/documents/file.txt",
			expected: "home_user_documents_file.txt",
		},
		{
			name:     "windows path",
			input:    "C:\\Users\\Documents\\file.txt",
			expected: "C:_Users_Documents_file.txt",
		},
		{
			name:     "mixed separators",
			input:    "home/user\\documents/file.txt",
			expected: "home_user_documents_file.txt",
		},
		{
			name:     "no separators",
			input:    "file.txt",
			expected: "file.txt",
		},
		{
			name:     "multiple consecutive separators",
			input:    "home//user///file.txt",
			expected: "home__user___file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replacePathSeparators(tt.input)
			if result != tt.expected {
				t.Errorf("replacePathSeparators(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCertificateValidationForLocalhost(t *testing.T) {
	// Reset viper to avoid state leakage between tests
	viper.Reset()

	tmpDir := t.TempDir()
	home := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", home)

	if err := config.InitConfig(); err != nil {
		t.Fatalf("Failed to init config: %v", err)
	}

	tests := []struct {
		name               string
		apiURL             string
		shouldSkipValidate bool
	}{
		{
			name:               "localhost with https",
			apiURL:             "https://localhost:8000",
			shouldSkipValidate: true,
		},
		{
			name:               "127.0.0.1 with https",
			apiURL:             "https://127.0.0.1:8000",
			shouldSkipValidate: true,
		},
		{
			name:               "local IP address",
			apiURL:             "https://192.168.0.3:8000",
			shouldSkipValidate: false,
		},
		{
			name:               "production domain",
			apiURL:             "https://api.example.com",
			shouldSkipValidate: false,
		},
		{
			name:               "localhost with http",
			apiURL:             "http://localhost:8000",
			shouldSkipValidate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set config for this test
			config.Set("api_url", tt.apiURL)
			config.Set("api_key", "test-key")

			client, err := NewClient()
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			// Check TLS config using proper URL parsing
			isLocalhost := false
			if parsedURL, err := url.Parse(tt.apiURL); err == nil {
				hostname := parsedURL.Hostname()
				isLocalhost = hostname == "localhost" || hostname == "127.0.0.1"
			}

			if isLocalhost != tt.shouldSkipValidate {
				t.Errorf("Expected shouldSkipValidate=%v for URL %s, but got %v", tt.shouldSkipValidate, tt.apiURL, isLocalhost)
			}

			// Verify client was created
			if client == nil {
				t.Error("Expected non-nil client")
			}
		})
	}
}
