package crypt

import (
	"log"

	"github.com/letheanVPN/desktop/services/crypt/lib/openpgp"
	"github.com/letheanVPN/desktop/services/filesystem"
)

// CryptService provides cryptographic functions
type CryptService struct{}

// NewCryptService creates a new CryptService
func NewCryptService() *CryptService {
	return &CryptService{}
}

// Startup is called when the app starts. This is a good place to do initialization
func (s *CryptService) Startup() {
	// This should be handled more gracefully, but for now, we'll panic on error.
	if err := filesystem.EnsureDir("config"); err != nil {
		log.Fatalf("Failed to create config directory: %v", err)
	}
	if err := filesystem.EnsureDir("data"); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}
	if err := filesystem.EnsureDir("workspaces"); err != nil {
		log.Fatalf("Failed to create workspaces directory: %v", err)
	}

	// Create server key pair if it doesn't exist
	if !filesystem.IsFile("users/server.lthn.pub") {
		log.Println("Creating server key pair...")
		if err := openpgp.CreateServerKeyPair(); err != nil {
			log.Fatalf("Failed to create server key pair: %v", err)
		}
		log.Println("Server key pair created.")
	}
}
