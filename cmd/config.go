/*
Copyright © 2026 T. Vicente <thiagoaureliovicente@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vcnt/sfs-cli/internal/config"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage sfs-cli configuration",
	Long: `Manage configuration for sfs-cli including API URL and API key.

Configuration is stored in ~/.sfs-cli.yaml

Available commands:
  set <key> <value>  Set a configuration value
  get <key>          Get a configuration value
  list               List all configuration`,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value. Available keys:
  api_url  - The base URL of the SFS API (default: https://localhost)
  api_key  - Your API key for authentication`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		if err := config.Set(key, value); err != nil {
			return fmt.Errorf("failed to set config: %w", err)
		}

		fmt.Printf("Configuration updated: %s = %s\n", key, value)
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := config.GetValue(key)

		if value == "" {
			fmt.Printf("%s is not set\n", key)
		} else {
			fmt.Printf("%s = %s\n", key, value)
		}
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		all := config.GetAll()

		if len(all) == 0 {
			fmt.Println("No configuration set")
			return nil
		}

		fmt.Println("Current configuration:")
		for key, value := range all {
			fmt.Printf("  %s = %v\n", key, value)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
}
