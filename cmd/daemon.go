/*
Copyright © 2026 T. Vicente <thiagoaureliovicente@gmail.com>

*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/vcnt/sfs-cli/internal/daemon"
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Manage the SFS background daemon service",
	Long: `The daemon command manages the SFS background service that enables
automatic file watching and synchronization.

Available subcommands:
  create  - Create the systemd service file
  enable  - Enable daemon to start automatically on boot
  disable - Disable automatic startup
  start   - Manually start the daemon
  stop    - Stop the daemon
  restart - Restart the daemon
  status  - Check daemon status

Note: This command is only supported on Linux systems.`,
}

const (
	serviceName = "sfs-daemon"
)

// checkLinux verifies that the OS is Linux
func checkLinux() error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("daemon command is only supported on Linux (current OS: %s)", runtime.GOOS)
	}
	return nil
}

// getServiceFilePath returns the path to the user systemd service file
func getServiceFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}

	serviceDir := filepath.Join(homeDir, ".config", "systemd", "user")
	return filepath.Join(serviceDir, serviceName+".service"), nil
}

// getServiceDir returns the user systemd service directory
func getServiceDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}

	return filepath.Join(homeDir, ".config", "systemd", "user"), nil
}

// getExecutablePath returns the absolute path to the current executable
func getExecutablePath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %v", err)
	}
	return filepath.Abs(execPath)
}

// generateServiceFile returns the systemd service file content
func generateServiceFile(execPath string) string {
	return fmt.Sprintf(`[Unit]
Description=Semantic File Search Daemon
After=network.target

[Service]
Type=simple
ExecStart=%s daemon run
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=default.target
`, execPath)
}

// runSystemctl executes a systemctl command with --user flag
func runSystemctl(args ...string) error {
	// Prepend --user to all systemctl commands
	cmdArgs := append([]string{"--user"}, args...)
	cmd := exec.Command("systemctl", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

var daemonCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create the systemd user service file",
	Long:  `Creates the systemd user service file at ~/.config/systemd/user/sfs-daemon.service.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkLinux(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		execPath, err := getExecutablePath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		serviceDir, err := getServiceDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Create the service directory if it doesn't exist
		if err := os.MkdirAll(serviceDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create service directory: %v\n", err)
			os.Exit(1)
		}

		serviceFilePath, err := getServiceFilePath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		serviceContent := generateServiceFile(execPath)

		if err := os.WriteFile(serviceFilePath, []byte(serviceContent), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create service file: %v\n", err)
			os.Exit(1)
		}

		if err := runSystemctl("daemon-reload"); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to reload systemd: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Service file created at %s\n", serviceFilePath)
		fmt.Println("Run 'sfs daemon enable' to enable autostart on boot")
	},
}

var daemonEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable daemon to start automatically on boot",
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkLinux(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := runSystemctl("enable", serviceName); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to enable daemon: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Daemon enabled successfully")
	},
}

var daemonDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable automatic startup",
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkLinux(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := runSystemctl("disable", serviceName); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to disable daemon: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Daemon disabled successfully")
	},
}

var daemonStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Manually start the daemon",
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkLinux(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := runSystemctl("start", serviceName); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to start daemon: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Daemon started successfully")
	},
}

var daemonStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the daemon",
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkLinux(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := runSystemctl("stop", serviceName); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to stop daemon: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Daemon stopped successfully")
	},
}

var daemonRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the daemon",
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkLinux(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := runSystemctl("restart", serviceName); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to restart daemon: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Daemon restarted successfully")
	},
}

var daemonStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check daemon status",
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkLinux(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := runSystemctl("status", serviceName); err != nil {
			os.Exit(1)
		}
	},
}

var daemonRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the daemon (used by systemd)",
	Long:  `This command is called by systemd to run the daemon. Do not call this directly.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := daemon.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Daemon error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)

	// Add subcommands
	daemonCmd.AddCommand(daemonCreateCmd)
	daemonCmd.AddCommand(daemonEnableCmd)
	daemonCmd.AddCommand(daemonDisableCmd)
	daemonCmd.AddCommand(daemonStartCmd)
	daemonCmd.AddCommand(daemonStopCmd)
	daemonCmd.AddCommand(daemonRestartCmd)
	daemonCmd.AddCommand(daemonStatusCmd)
	daemonCmd.AddCommand(daemonRunCmd)
}
