package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	// Use temp dir for test
	tmpDir := t.TempDir()
	home := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", home)

	// Initialize config
	if err := InitConfig(); err != nil {
		t.Fatalf("Failed to init config: %v", err)
	}

	// Test Set
	testKey := "test_key"
	testValue := "test_value"

	if err := Set(testKey, testValue); err != nil {
		t.Fatalf("Failed to set config: %v", err)
	}

	// Test Get
	value := GetValue(testKey)
	if value != testValue {
		t.Errorf("Expected %s, got %s", testValue, value)
	}

	// Verify file was created
	configFile := filepath.Join(tmpDir, ConfigDirName, ConfigFileName+"."+ConfigFileType)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}
}

func TestGetDefault(t *testing.T) {
	// Use temp dir for test
	tmpDir := t.TempDir()
	home := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", home)

	// Initialize config
	if err := InitConfig(); err != nil {
		t.Fatalf("Failed to init config: %v", err)
	}

	// Test default values
	cfg, err := Get()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	if cfg.APIURL != "https://localhost" {
		t.Errorf("Expected default API URL 'https://localhost', got '%s'", cfg.APIURL)
	}

	if cfg.APIKey != "" {
		t.Errorf("Expected empty API key, got '%s'", cfg.APIKey)
	}
}

func TestGetAll(t *testing.T) {
	// Use temp dir for test
	tmpDir := t.TempDir()
	home := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", home)

	// Initialize config
	if err := InitConfig(); err != nil {
		t.Fatalf("Failed to init config: %v", err)
	}

	// Set some values
	Set("key1", "value1")
	Set("key2", "value2")

	// Get all
	all := GetAll()

	if len(all) == 0 {
		t.Error("Expected non-empty config map")
	}

	if all["key1"] != "value1" {
		t.Errorf("Expected key1='value1', got '%v'", all["key1"])
	}
}

func TestConfigFilePermissions(t *testing.T) {
	// Use temp dir for test
	tmpDir := t.TempDir()
	home := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", home)

	// Initialize config
	if err := InitConfig(); err != nil {
		t.Fatalf("Failed to init config: %v", err)
	}

	// Set a value to create the config file
	if err := Set("api_key", "secret-key"); err != nil {
		t.Fatalf("Failed to set config: %v", err)
	}

	// Check file permissions
	configFile := filepath.Join(tmpDir, ConfigDirName, ConfigFileName+"."+ConfigFileType)
	info, err := os.Stat(configFile)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}

	// Check that permissions are 0600 (owner read/write only)
	mode := info.Mode().Perm()
	if mode != 0600 {
		t.Errorf("Expected file permissions 0600, got %04o", mode)
	}
}
