package sftp

import (
	"github.com/pkg/sftp"
)

// Medium implements the filesystem.Medium interface for the SFTP protocol.
type Medium struct {
	client *sftp.Client
}

// ConnectionConfig holds the necessary details to connect to an SFTP server.
type ConnectionConfig struct {
	Host     string
	Port     string
	User     string
	Password string // For password-based auth
	KeyFile  string // Path to a private key for key-based auth
}
