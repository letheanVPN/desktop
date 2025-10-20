package local

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// New creates a new instance of the local storage medium.
// It requires a root path to sandbox all file operations.
func New(rootPath string) (*Medium, error) {
	if err := os.MkdirAll(rootPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("could not create root directory at %s: %w", rootPath, err)
	}
	return &Medium{root: rootPath}, nil
}

// path returns a full, safe path within the medium's root.
func (m *Medium) path(subpath string) (string, error) {
	if strings.Contains(subpath, "..") {
		return "", fmt.Errorf("path traversal attempt detected")
	}
	return filepath.Join(m.root, subpath), nil
}

// Read retrieves the content of a file from the local disk.
func (m *Medium) Read(path string) (string, error) {
	safePath, err := m.path(path)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(safePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Write saves the given content to a file on the local disk.
func (m *Medium) Write(path, content string) error {
	safePath, err := m.path(path)
	if err != nil {
		return err
	}
	dir := filepath.Dir(safePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(safePath, []byte(content), 0644)
}

// EnsureDir makes sure a directory exists on the local disk.
func (m *Medium) EnsureDir(path string) error {
	safePath, err := m.path(path)
	if err != nil {
		return err
	}
	return os.MkdirAll(safePath, os.ModePerm)
}

// IsFile checks if a path exists and is a regular file on the local disk.
func (m *Medium) IsFile(path string) bool {
	safePath, err := m.path(path)
	if err != nil {
		return false
	}
	info, err := os.Stat(safePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
