/*
Copyright © 2026 T. Vicente<thiagoaureliovicente@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ThiagoAVicente/sfs-cli/internal/api"
)

var updateFlag bool

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload <file>",
	Short: "Upload a file to the SFS API for indexing",
	Long: `Upload a file to the SFS API for semantic indexing.

The file will be processed and indexed, making it searchable via semantic queries.

Examples:
  sfs upload document.pdf
  sfs upload --update existing_file.txt    # Update existing file`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		client, err := api.NewClient()
		if err != nil {
			return err
		}

		_, err = client.UploadFile(filePath, updateFlag)
		return err
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().BoolVarP(&updateFlag, "update", "u", false, "Update existing file")
}
