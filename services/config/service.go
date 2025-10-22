package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/letheanVPN/desktop/services/core/i18n"
)

const appName = "lethean"
const configFileName = "config.json"

// ErrSetupRequired is returned by ServiceStartup if config.json is missing.
var ErrSetupRequired = errors.New("setup required: config.json not found")

// Service provides access to the application's configuration.
type Service struct {
	config *Config
}

// storableConfig defines the structure of the data that is actually written to config.json.
// This prevents dynamic paths from being written to the config file.
type storableConfig struct {
	DefaultRoute string   `json:"defaultRoute,omitempty"`
	Features     []string `json:"features,omitempty"`
	Language     string   `json:"language,omitempty"`
}

// NewService creates and initializes a new configuration service.
// It loads an existing configuration or creates a default one if not found.
func NewService() (*Service, error) {
	// 1. Determine the config directory path to check for an existing file.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not resolve user home directory: %w", err)
	}
	userHomeDir := filepath.Join(homeDir, appName)
	configDir := filepath.Join(userHomeDir, "config")
	configPath := filepath.Join(configDir, configFileName)

	var cfg *Config
	configNeedsSaving := false

	// 2. Check if the config file exists.
	if _, err := os.Stat(configPath); err == nil {
		// --- Config file EXISTS ---

		// First, get the base config with all the dynamic paths and directory structures.
		cfg, err = newDefaultConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to create base config structure: %w", err)
		}

		// Now, load the storable values from the existing file.
		fileData, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read existing config file at %s: %w", configPath, err)
		}

		loadedStorableConfig := &storableConfig{}
		if err := json.Unmarshal(fileData, loadedStorableConfig); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to unmarshal config.json at %s, using defaults: %v\n", configPath, err)
		} else {
			// Apply loaded storable values over the defaults.
			if loadedStorableConfig.DefaultRoute != "" {
				cfg.DefaultRoute = loadedStorableConfig.DefaultRoute
			}
			if loadedStorableConfig.Features != nil {
				cfg.Features = loadedStorableConfig.Features
			}
			if loadedStorableConfig.Language != "" {
				cfg.Language = loadedStorableConfig.Language
			}
		}
	} else if errors.Is(err, os.ErrNotExist) {
		// --- Config file DOES NOT EXIST ---
		configNeedsSaving = true

		// Create a fresh default config. This sets up paths and a default "en" language.
		cfg, err = newDefaultConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}

		// Now, perform the "expensive" operation of detecting the language.
		if detectedLang, err := i18n.DetectLanguage(); err == nil && detectedLang != "" {
			cfg.Language = detectedLang
		}

	} else {
		// Another error occurred (e.g., permissions).
		return nil, fmt.Errorf("failed to check for config file at %s: %w", configPath, err)
	}

	service := &Service{config: cfg}

	// If the config file didn't exist, save the newly generated one.
	if configNeedsSaving {
		if err := service.Save(); err != nil {
			return nil, fmt.Errorf("failed to save initial config: %w", err)
		}
	}

	return service, nil
}

// newDefaultConfig creates a default configuration with resolved paths and ensures directories exist.
// It sets a hardcoded default language "en".
func newDefaultConfig() (*Config, error) {
	if strings.Contains(appName, "..") || strings.Contains(appName, string(filepath.Separator)) {
		return nil, fmt.Errorf("invalid app name '%s': contains path traversal characters", appName)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not resolve user home directory: %w", err)
	}
	userHomeDir := filepath.Join(homeDir, appName)

	rootDir, err := xdg.DataFile(appName)
	if err != nil {
		return nil, fmt.Errorf("could not resolve data directory: %w", err)
	}

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
		DefaultRoute:  "/",
		Features:      []string{},
		Language:      "en", // Hardcoded default, will be overridden if loaded or detected
	}

	dirs := []string{cfg.RootDir, cfg.ConfigDir, cfg.DataDir, cfg.CacheDir, cfg.WorkspacesDir, cfg.UserHomeDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("could not create directory %s: %w", dir, err)
		}
	}

	return cfg, nil
}

// Get returns the loaded configuration.
func (s *Service) Get() *Config {
	return s.config
}

// Load reads the configuration from config.json.
// This method is now primarily used internally by NewService to apply loaded values.
func (s *Service) Load() error {
	configPath := filepath.Join(s.config.ConfigDir, configFileName)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	stored := &storableConfig{}
	if err := json.Unmarshal(data, stored); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Update only the fields that are user-configurable
	if stored.DefaultRoute != "" {
		s.config.DefaultRoute = stored.DefaultRoute
	}
	if stored.Features != nil {
		s.config.Features = stored.Features
	}
	// Language is handled explicitly in NewService to manage detection vs loaded
	if stored.Language != "" {
		s.config.Language = stored.Language
	}

	return nil
}

// Save writes the current configuration to config.json.
func (s *Service) Save() error {
	configPath := filepath.Join(s.config.ConfigDir, configFileName)

	// Only save the storable parts of the config
	stored := &storableConfig{
		DefaultRoute: s.config.DefaultRoute,
		Features:     s.config.Features,
		Language:     s.config.Language,
	}

	data, err := json.MarshalIndent(stored, "", "  ")
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
	if s.IsFeatureEnabled(feature) {
		return nil
	}
	s.config.Features = append(s.config.Features, feature)
	if err := s.Save(); err != nil {
		return fmt.Errorf("failed to save config after enabling feature %s: %w", feature, err)
	}
	return nil
}
