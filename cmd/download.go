/*
Copyright © 2026 T. Vicente<thiagoaureliovicente@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ThiagoAVicente/sfs-cli/internal/api"
)

var outputPath string

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download <filename> [output]",
	Short: "Download a file from the SFS system",
	Long: `Download a file that has been stored in the SFS system.

If output path is not specified, the file will be downloaded with its original name.

Examples:
  sfs download document.pdf
  sfs download home_user_docs_notes.txt ./notes.txt
  sfs download file.txt --output ./downloaded.txt`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		fileName := args[0]

		// Determine output path
		dest := outputPath
		if dest == "" && len(args) > 1 {
			dest = args[1]
		}
		if dest == "" {
			dest = fileName
		}

		client, err := api.NewClient()
		if err != nil {
			return err
		}

		if err := client.DownloadFile(fileName, dest); err != nil {
			return err
		}

		fmt.Printf("File downloaded: %s -> %s\n", fileName, dest)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path")
}
