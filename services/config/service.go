package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
)

const appName = "lethean"

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

	// Now define all other paths relative to our application root.
	cfg := &Config{
		UserHomeDir:   userHomeDir,
		RootDir:       rootDir,
		CacheDir:      cacheDir,
		ConfigDir:     filepath.Join(userHomeDir, "config"),
		DataDir:       filepath.Join(userHomeDir, "data"),
		WorkspacesDir: filepath.Join(userHomeDir, "workspaces"),
	}

	// Ensure all base directories exist using the standard os library.
	// This makes the config service self-sufficient.
	dirs := []string{cfg.RootDir, cfg.ConfigDir, cfg.DataDir, cfg.CacheDir, cfg.WorkspacesDir, cfg.UserHomeDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("could not create directory %s: %w", dir, err)
		}
	}

	return &Service{config: cfg}, nil
}

// Get returns the loaded configuration.
func (s *Service) Get() *Config {
	return s.config
}
