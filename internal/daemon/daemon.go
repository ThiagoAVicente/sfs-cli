/*
Copyright © 2026 T. Vicente <thiagoaureliovicente@gmail.com>

*/
package daemon

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/vcnt/sfs-cli/internal/config"
)

func Run() error {
	log.Println("SFS daemon starting...")

	// Create config watcher
	configWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer configWatcher.Close()

	// Get config directory and file path
	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}
	configPath, err := config.GetConfigPath()
	if err != nil {
		return err
	}
	configFileName := filepath.Base(configPath)

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Printf("Warning: Could not create config directory: %v", err)
	}

	// Watch the config directory
	err = configWatcher.Add(configDir)
	if err != nil {
		log.Printf("Warning: Could not watch config directory %s: %v", configDir, err)
	} else {
		log.Printf("Watching config file: %s", configPath)
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Load initial config
	if err := config.InitConfig(); err != nil {
		log.Printf("Warning: Failed to load config: %v", err)
	}

	// TODO: Initialize file watcher here (you'll implement this)

	log.Println("Daemon is running. Press Ctrl+C to stop.")

	// Main event loop
	for {
		select {
		case event, ok := <-configWatcher.Events:
			if !ok {
				return nil
			}

			if filepath.Base(event.Name) != configFileName {
				continue
			}

			// Handle config file changes
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				log.Printf("Config file changed: %s", event.Name)

				// Reload config
				if err := config.InitConfig(); err != nil {
					log.Printf("Error reloading config: %v", err)
				} else {
					log.Println("Config reloaded successfully")
					// TODO: Update file watcher with new watched directories
				}
			}

		case err, ok := <-configWatcher.Errors:
			if !ok {
				return nil
			}
			log.Printf("Config watcher error: %v", err)

		case sig := <-sigChan:
			log.Printf("Received signal %v, shutting down...", sig)
			return nil
		}
	}
}
