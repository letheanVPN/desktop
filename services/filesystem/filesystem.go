package filesystem

import (
	"log"

	"github.com/letheanVPN/desktop/services/core/config"
	"github.com/letheanVPN/desktop/services/filesystem/local"
)

// Medium defines the standard interface for a storage backend.
// This allows for different implementations (e.g., local disk, S3, SFTP)
// to be used interchangeably.
type Medium interface {
	// Read retrieves the content of a file as a string.
	Read(path string) (string, error)

	// Write saves the given content to a file, overwriting it if it exists.
	Write(path, content string) error

	// EnsureDir makes sure a directory exists, creating it if necessary.
	EnsureDir(path string) error

	// IsFile checks if a path exists and is a regular file.
	IsFile(path string) bool
}

// Pre-initialized, sandboxed medium for the local filesystem.
var Local Medium

// init runs once when the package is first used, setting up the Local medium.
func init() {
	configService, err := config.NewService()
	if err != nil {
		log.Fatalf("Fatal: Filesystem could not be initialized: %v", err)
	}
	cfg := configService.Get()

	Local, err = local.New(cfg.RootDir)
	if err != nil {
		log.Fatalf("Fatal: Local filesystem medium could not be created: %v", err)
	}
}
