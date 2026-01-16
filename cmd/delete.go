/*
Copyright © 2026 T. Vicente<thiagoaureliovicente@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ThiagoAVicente/sfs-cli/internal/api"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <filename>",
	Short: "Delete an indexed file",
	Long: `Delete a file from the SFS system.

This will remove both the file and its index data.

Examples:
  sfs delete document.pdf
  sfs delete home_user_docs_notes.txt`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fileName := args[0]

		client, err := api.NewClient()
		if err != nil {
			return err
		}

		result, err := client.DeleteFile(fileName)
		if err != nil {
			return err
		}

		fmt.Printf("File deleted: %s\n", fileName)
		fmt.Printf("Job ID: %s\n", result.JobID)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
