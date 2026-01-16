/*
Copyright © 2026 T. Vicente<thiagoaureliovicente@gmail.com>

*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vcnt/sfs-cli/internal/config"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Manage directories to watch for automatic syncing",
	Long: `Manage the list of directories that the daemon watches for automatic file syncing.

When you add directories to watch, the daemon will automatically upload any changes
to the SFS API.`,
}

var watchAddCmd = &cobra.Command{
	Use:   "add <directory>",
	Short: "Add a directory to watch",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]

		// Convert to absolute path
		absDir, err := filepath.Abs(dir)
		if err != nil {
			return fmt.Errorf("failed to resolve path: %w", err)
		}

		// Check if directory exists
		if info, err := os.Stat(absDir); err != nil {
			return fmt.Errorf("directory does not exist: %s", absDir)
		} else if !info.IsDir() {
			return fmt.Errorf("path is not a directory: %s", absDir)
		}

		// Load config
		if err := config.InitConfig(); err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Get current watch dirs
		watchDirs := config.GetWatchDirs()

		if slices.Contains(watchDirs, absDir) {
			fmt.Printf("Directory already being watched: %s\n", absDir)
			return nil
		}

		// Add to list
		watchDirs = append(watchDirs, absDir)
		viper.Set("watch_dirs", watchDirs)

		// Save config
		if err := config.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Added to watch list: %s\n", absDir)
		return nil
	},
}

var watchRemoveCmd = &cobra.Command{
	Use:   "remove <directory>",
	Short: "Remove a directory from watch list",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]

		// Convert to absolute path
		absDir, err := filepath.Abs(dir)
		if err != nil {
			return fmt.Errorf("failed to resolve path: %w", err)
		}

		// Load config
		if err := config.InitConfig(); err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Get current watch dirs
		watchDirs := config.GetWatchDirs()

		// Find and remove
		found := false
		newWatchDirs := []string{}
		for _, d := range watchDirs {
			if d == absDir {
				found = true
			} else {
				newWatchDirs = append(newWatchDirs, d)
			}
		}

		if !found {
			return fmt.Errorf("directory not in watch list: %s", absDir)
		}

		// Update config
		viper.Set("watch_dirs", newWatchDirs)

		// Save config
		if err := config.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Removed from watch list: %s\n", absDir)
		return nil
	},
}

var watchListCmd = &cobra.Command{
	Use:   "list",
	Short: "List watched directories",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load config
		if err := config.InitConfig(); err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Get watch dirs
		watchDirs := config.GetWatchDirs()

		if len(watchDirs) == 0 {
			fmt.Println("No directories being watched")
			fmt.Println("Add directories with: sfs watch add <directory>")
			return nil
		}

		fmt.Println("Watched directories:")
		for _, dir := range watchDirs {
			fmt.Printf("  %s\n", dir)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
	watchCmd.AddCommand(watchAddCmd)
	watchCmd.AddCommand(watchRemoveCmd)
	watchCmd.AddCommand(watchListCmd)
}
