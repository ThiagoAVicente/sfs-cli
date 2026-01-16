/*
Copyright © 2026 T. Vicente <thiagoaureliovicente@gmail.com>
*/
package daemon

import (
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/vcnt/sfs-cli/internal/api"
	"github.com/vcnt/sfs-cli/internal/config"
)

const debounceDelay = 500 * time.Millisecond

var (
	debounceTimers   = make(map[string]*time.Timer)
	debounceMutex    sync.Mutex
)

func ensure(err error, msg string, stopOnErr bool) {
	if err != nil {
		log.Printf("%s: %v", msg, err)
		if stopOnErr {
			os.Exit(1)
		}
	}
}

// receives a list of directories/files and adds to watcher
func create_watcher(dirs []string) *fsnotify.Watcher {
	fileWatcher, err := fsnotify.NewWatcher()
	ensure(err, "Failed to create file watcher", true)

	for _, dir := range dirs {
		// Convert to absolute path
		absDir, err := filepath.Abs(dir)
		if err != nil {
			log.Printf("Warning: Could not resolve path %s: %v", dir, err)
			continue
		}

		// Walk recursively to add all subdirectories
		filepath.WalkDir(absDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Printf("Error walking %s: %v", path, err)
				return nil
			}
			if d.IsDir() {
				if err := fileWatcher.Add(path); err != nil {
					log.Printf("Warning: Could not watch %s: %v", path, err)
				} else {
					log.Printf("Watching: %s", path)
				}
			}
			return nil
		})
	}
	return fileWatcher
}

func Run() error {
	log.Println("SFS daemon starting...")

	// Create config watcher
	configWatcher, err := fsnotify.NewWatcher()
	ensure(err, "Failed to create config watcher", true)
	defer configWatcher.Close()

	// Get config directory and file path
	configDir, err := config.GetConfigDir()
	ensure(err, "Failed to get config directory", true)

	configPath, err := config.GetConfigPath()
	ensure(err, "Failed to get config path", true)

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

	// crearte file watcher
	fileWatcher := create_watcher(config.GetWatchDirs())
	defer fileWatcher.Close()

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
					fileWatcher.Close()
					fileWatcher = create_watcher(config.GetWatchDirs())
				}
			}

		case event, ok := <-fileWatcher.Events:
			if !ok {
				return nil
			}

			// Handle file changes
			if event.Has(fsnotify.Write) {
				// Skip backup/temp files
				if strings.HasSuffix(event.Name, "~") || strings.HasSuffix(event.Name, ".swp") {
					continue
				}

				// Skip directories
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					continue
				}

				log.Printf("File changed: %s (debouncing...)", event.Name)

				// Debounce: cancel existing timer and set new one
				debounceMutex.Lock()
				if timer, exists := debounceTimers[event.Name]; exists {
					timer.Stop()
				}

				debounceTimers[event.Name] = time.AfterFunc(debounceDelay, func() {
					// Upload the file after delay
					cli, err := api.NewClient()
					if err != nil {
						log.Printf("Failed to create client: %v", err)
						return
					}

					if _, err := cli.UploadFile(event.Name, true); err != nil {
						log.Printf("Failed to upload file %s: %v", event.Name, err)
					} else {
						log.Printf("Uploaded file: %s", event.Name)
					}

					// Clean up timer
					debounceMutex.Lock()
					delete(debounceTimers, event.Name)
					debounceMutex.Unlock()
				})
				debounceMutex.Unlock()
			}

		case err, ok := <-fileWatcher.Errors:
			if !ok {
				return nil
			}
			log.Printf("File watcher error: %v", err)

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
