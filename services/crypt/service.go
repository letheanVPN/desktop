package crypt

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/letheanVPN/desktop/services/core/config"
	"github.com/letheanVPN/desktop/services/crypt/lib/openpgp"
	"github.com/letheanVPN/desktop/services/filesystem"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// createServerKeyPair is a package-level variable that can be swapped for testing.
var createServerKeyPair = openpgp.CreateServerKeyPair

// NewService creates a new crypt.Service, accepting a config service instance.
func NewService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
	}
}

// ServiceStartup Startup is called when the app starts. It handles one-time cryptographic setup.
func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	// Define the directory for server keys based on the central config.
	serverKeysDir := filepath.Join(s.config.DataDir, "server_keys")
	if err := filesystem.EnsureDir(filesystem.Local, serverKeysDir); err != nil {
		return fmt.Errorf("failed to create server keys directory: %w", err)
	}

	// Check for server key pair using the configured path.
	serverKeyPath := filepath.Join(serverKeysDir, "server.lthn.pub")
	if !filesystem.IsFile(filesystem.Local, serverKeyPath) {
		log.Println("Creating server key pair...")
		if err := createServerKeyPair(serverKeysDir); err != nil {
			return fmt.Errorf("failed to create server key pair: %w", err)
		}
		log.Println("Server key pair created.")
	}
	return nil
}
