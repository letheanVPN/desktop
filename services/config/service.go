package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
)

const appName = "lethean"
const configFileName = "config.json"

// ErrSetupRequired is returned by ServiceStartup if config.json is missing.
var ErrSetupRequired = errors.New("setup required: config.json not found")

// Service provides access to the application's configuration.
type Service struct {
	config *Config
}

// NewService creates and initializes a new configuration service.
// It resolves OS-specific paths and ensures they exist.
func NewService() (*Service, error) {
	// Validate appName for path traversal attempts
	if strings.Contains(appName, "..") || strings.Contains(appName, string(filepath.Separator)) {
		return nil, fmt.Errorf("invalid app name '%s': contains path traversal characters", appName)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not resolve user home directory: %w", err)
	}
	// This is the root for user-configurable data, like workspaces, configs, etc.
	userHomeDir := filepath.Join(homeDir, appName)

	// Use XDG for the application's own files (e.g., binaries, assets).
	// On macOS, this resolves to ~/Library/Application Support/lethean
	// On Linux, this resolves to ~/.local/share/lethean
	rootDir, err := xdg.DataFile(appName)
	if err != nil {
		return nil, fmt.Errorf("could not resolve data directory: %w", err)
	}

	// Per XDG specs, cache should be in a different location.
	cacheDir, err := xdg.CacheFile(appName)
	if err != nil {
		return nil, fmt.Errorf("could not resolve cache directory: %w", err)
	}

	cfg := &Config{
		UserHomeDir:   userHomeDir,
		RootDir:       rootDir,
		CacheDir:      cacheDir,
		ConfigDir:     filepath.Join(userHomeDir, "config"),
		DataDir:       filepath.Join(userHomeDir, "data"),
		WorkspacesDir: filepath.Join(userHomeDir, "workspaces"),
		DefaultRoute:  "/",        // Set default route here
		Features:      []string{}, // Initialize empty features
	}

	// Ensure all base directories exist using the standard os library.
	// This makes the config service self-sufficient.
	dirs := []string{cfg.RootDir, cfg.ConfigDir, cfg.DataDir, cfg.CacheDir, cfg.WorkspacesDir, cfg.UserHomeDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("could not create directory %s: %w", dir, err)
		}
	}

	service := &Service{config: cfg}

	// Attempt to load existing config.json
	if err := service.Load(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// If config.json doesn't exist, save the default config
			if err := service.Save(); err != nil {
				return nil, fmt.Errorf("failed to save initial config: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	}

	return service, nil
}

// Get returns the loaded configuration.
func (s *Service) Get() *Config {
	return s.config
}

// Load reads the configuration from config.json.
func (s *Service) Load() error {
	configPath := filepath.Join(s.config.ConfigDir, configFileName)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Create a temporary config to unmarshal into, preserving existing paths
	tempConfig := &Config{}
	if err := json.Unmarshal(data, tempConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Update only the fields that can be loaded from the file
	s.config.DefaultRoute = tempConfig.DefaultRoute
	s.config.Features = tempConfig.Features

	return nil
}

// Save writes the current configuration to config.json.
func (s *Service) Save() error {
	configPath := filepath.Join(s.config.ConfigDir, configFileName)
	data, err := json.MarshalIndent(s.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

// IsFeatureEnabled checks if a given feature is enabled in the configuration.
func (s *Service) IsFeatureEnabled(feature string) bool {
	for _, f := range s.config.Features {
		if f == feature {
			return true
		}
	}
	return false
}

// EnableFeature adds a feature to the list of enabled features and saves the config.
func (s *Service) EnableFeature(feature string) error {
	// Check if the feature is already enabled
	if s.IsFeatureEnabled(feature) {
		return nil // Feature already enabled, nothing to do
	}

	// Add the feature
	s.config.Features = append(s.config.Features, feature)

	// Save the updated config
	if err := s.Save(); err != nil {
		return fmt.Errorf("failed to save config after enabling feature %s: %w", feature, err)
	}

	return nil
}
