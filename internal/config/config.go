package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	ConfigFileName = ".sfs-cli"
	ConfigFileType = "yaml"
)

// Config holds the application configuration
type Config struct {
	APIURL string `mapstructure:"api_url"`
	APIKey string `mapstructure:"api_key"`
}

// InitConfig initializes viper configuration
func InitConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Set config file location
	viper.AddConfigPath(home)
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileType)

	// Set defaults
	viper.SetDefault("api_url", "https://localhost")
	viper.SetDefault("api_key", "")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		// Config file not found is okay, we'll create it on first set
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config: %w", err)
		}
	}

	return nil
}

// Get returns the current configuration
func Get() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &cfg, nil
}

// Set sets a configuration value
func Set(key, value string) error {
	viper.Set(key, value)
	return Save()
}

// GetValue gets a single configuration value
func GetValue(key string) string {
	return viper.GetString(key)
}

// Save saves the current configuration to disk
func Save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(home, ConfigFileName+"."+ConfigFileType)

	if err := viper.WriteConfigAs(configPath); err != nil {
		// If file doesn't exist, SafeWriteConfig creates it
		if err := viper.SafeWriteConfig(); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
	}

	// Set secure permissions (owner read/write only)
	if err := os.Chmod(configPath, 0600); err != nil {
		return fmt.Errorf("failed to set config permissions: %w", err)
	}

	return nil
}

// GetAll returns all configuration as a map
func GetAll() map[string]interface{} {
	return viper.AllSettings()
}
