/*
Copyright © 2026 T. Vicente <thiagoaureliovicente@gmail.com>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ThiagoAVicente/sfs-cli/internal/config"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sfs",
	Short: "CLI tool for Semantic File Search",
	Long: `sfs-cli is a command-line tool to interact with the SFS (Semantic File Search) API.

Use it to upload files, search semantically, manage indexed files, and more.

Configuration:
  Config file: ~/.config/sfs/config.yaml

  Required settings:
    api_url - Base URL of your SFS API
    api_key - Your API authentication key

Examples:
  sfs config set api_url https://api.example.com
  sfs config set api_key your-secret-key
  sfs upload /path/to/file.txt
  sfs search "find relevant documents"`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if err := config.InitConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize config: %v\n", err)
	}
}
