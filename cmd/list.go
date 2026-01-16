/*
Copyright © 2026 T. Vicente<thiagoaureliovicente@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ThiagoAVicente/sfs-cli/internal/api"
)

var prefixFilter string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all indexed files",
	Long: `List all files that have been indexed in the SFS system.

Optionally filter by filename prefix.

Examples:
  sfs list
  sfs list --prefix docs_
  sfs list --prefix home_user_`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.ListFiles(prefixFilter)
		if err != nil {
			return err
		}

		if len(result.Files) == 0 {
			fmt.Println("No files found")
			return nil
		}

		fmt.Printf("Found %d files:\n\n", result.Count)
		for _, file := range result.Files {
			fmt.Printf("  - %s\n", file)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&prefixFilter, "prefix", "p", "", "Filter files by prefix")
}
