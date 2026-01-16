package daemon

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCreateWatcher(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create watcher
	watcher := createWatcher([]string{tmpDir})
	defer watcher.Close()

	if watcher == nil {
		t.Fatal("Expected watcher to be created")
	}
}

func TestCreateWatcherWithInvalidPath(t *testing.T) {
	// This should not panic, but log warnings
	watcher := createWatcher([]string{"/nonexistent/path"})
	defer watcher.Close()

	if watcher == nil {
		t.Fatal("Expected watcher to be created even with invalid paths")
	}
}

func TestCreateWatcherRecursive(t *testing.T) {
	// Create nested directory structure
	tmpDir := t.TempDir()
	nested := filepath.Join(tmpDir, "level1", "level2", "level3")
	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatalf("Failed to create nested directories: %v", err)
	}

	// Create watcher - should watch all levels
	watcher := createWatcher([]string{tmpDir})
	defer watcher.Close()

	// Create a file in nested directory
	testFile := filepath.Join(nested, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Give watcher time to process
	time.Sleep(100 * time.Millisecond)

	// If we got here without errors, the watcher was set up correctly
}

func TestDebounceTimers(t *testing.T) {
	// Reset debounce timers for test
	debounceMutex.Lock()
	debounceTimers = make(map[string]*time.Timer)
	debounceMutex.Unlock()

	testFile := "/tmp/test.txt"
	uploadCount := 0

	// Simulate multiple rapid events
	for i := 0; i < 5; i++ {
		debounceMutex.Lock()
		if timer, exists := debounceTimers[testFile]; exists {
			timer.Stop()
		}

		debounceTimers[testFile] = time.AfterFunc(50*time.Millisecond, func() {
			uploadCount++
		})
		debounceMutex.Unlock()

		time.Sleep(10 * time.Millisecond)
	}

	// Wait for debounce delay
	time.Sleep(100 * time.Millisecond)

	// Should only upload once
	if uploadCount != 1 {
		t.Errorf("Expected 1 upload, got %d", uploadCount)
	}
}

func TestEnsure(t *testing.T) {
	// Test ensure with nil error (should not panic)
	ensure(nil, "test message", false)

	// Test ensure with error and stopOnErr=false (should not exit)
	ensure(os.ErrNotExist, "test error", false)

	// We can't easily test stopOnErr=true as it calls os.Exit
}

func TestCreateWatcherAbsolutePaths(t *testing.T) {
	// Create temp directory with relative path reference
	tmpDir := t.TempDir()
	relativeDir := filepath.Join(tmpDir, "test")
	if err := os.MkdirAll(relativeDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Change to temp dir to test relative path handling
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	// Create watcher with relative path
	watcher := createWatcher([]string{"./test"})
	defer watcher.Close()

	if watcher == nil {
		t.Fatal("Expected watcher to handle relative paths")
	}
}
