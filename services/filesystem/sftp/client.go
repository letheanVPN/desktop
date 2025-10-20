package sftp

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// New creates a new, connected instance of the SFTP storage medium.
func New(cfg ConnectionConfig) (*Medium, error) {
	var authMethods []ssh.AuthMethod

	if cfg.KeyFile != "" {
		key, err := os.ReadFile(cfg.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("unable to read private key: %w", err)
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("unable to parse private key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	} else if cfg.Password != "" {
		authMethods = append(authMethods, ssh.Password(cfg.Password))
	} else {
		return nil, fmt.Errorf("no authentication method provided (password or keyfile)")
	}

	sshConfig := &ssh.ClientConfig{
		User:            cfg.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // WARNING: In production, use a proper host key callback.
	}

	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial ssh: %w", err)
	}

	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		// Ensure the underlying ssh connection is closed on failure
		conn.Close()
		return nil, fmt.Errorf("failed to create sftp client: %w", err)
	}

	return &Medium{client: sftpClient}, nil
}

// Read retrieves the content of a file from the SFTP server.
func (m *Medium) Read(path string) (string, error) {
	file, err := m.client.Open(path)
	if err != nil {
		return "", fmt.Errorf("sftp: failed to open file %s: %w", path, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("sftp: failed to read file %s: %w", path, err)
	}

	return string(data), nil
}

// Write saves the given content to a file on the SFTP server.
func (m *Medium) Write(path, content string) error {
	// Ensure the remote directory exists first.
	dir := filepath.Dir(path)
	if err := m.EnsureDir(dir); err != nil {
		return err
	}

	file, err := m.client.Create(path)
	if err != nil {
		return fmt.Errorf("sftp: failed to create file %s: %w", path, err)
	}
	defer file.Close()

	if _, err := file.Write([]byte(content)); err != nil {
		return fmt.Errorf("sftp: failed to write to file %s: %w", path, err)
	}

	return nil
}

// EnsureDir makes sure a directory exists on the SFTP server.
func (m *Medium) EnsureDir(path string) error {
	// MkdirAll is idempotent, so it won't error if the path already exists.
	return m.client.MkdirAll(path)
}

// IsFile checks if a path exists and is a regular file on the SFTP server.
func (m *Medium) IsFile(path string) bool {
	info, err := m.client.Stat(path)
	if err != nil {
		// If the error is "not found", it's definitely not a file.
		// For any other error, we also conservatively say it's not a file.
		return false
	}
	// Return true only if it's not a directory.
	return !info.IsDir()
}
