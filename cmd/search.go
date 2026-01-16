/*
Copyright © 2026 T. Vicente<thiagoaureliovicente@gmail.com>

*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vcnt/sfs-cli/internal/api"
)

var (
	searchLimit     int
	scoreThreshold float64
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Perform semantic search on indexed files",
	Long: `Search your indexed files using semantic similarity.

The search uses vector embeddings to find relevant content based on meaning,
not just keyword matching.

Examples:
  sfs search "machine learning algorithms"
  sfs search "how to deploy applications" --limit 10
  sfs search "security best practices" --threshold 0.7`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")

		client, err := api.NewClient()
		if err != nil {
			return err
		}

		results, err := client.Search(query, searchLimit, scoreThreshold)
		if err != nil {
			return err
		}

		if len(results.Results) == 0 {
			fmt.Println("No results found")
			return nil
		}

		fmt.Printf("Found %d results:\n\n", len(results.Results))
		for i, result := range results.Results {
			fmt.Printf("[%d] Score: %.3f | File: %s\n", i+1, result.Score, result.Payload.FilePath)
			fmt.Printf("    Position: %d-%d | Chunk: %d\n", result.Payload.Start, result.Payload.End, result.Payload.ChunkIndex)
			fmt.Printf("    Text: %s\n\n", result.Payload.Text)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "l", 5, "Maximum number of results")
	searchCmd.Flags().Float64VarP(&scoreThreshold, "threshold", "t", 0.5, "Minimum similarity score (0.0-1.0)")
}
