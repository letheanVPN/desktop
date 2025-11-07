package workspace

import "github.com/Snider/Core/pkg/io"

// localMedium implements the Medium interface for the local disk.
type localMedium struct{}

// NewLocalMedium creates a new instance of the local storage medium.
func NewLocalMedium() io.Medium {
	return &localMedium{}
}

// FileGet reads a file from the local disk.
func (m *localMedium) FileGet(path string) (string, error) {
	return io.Read(io.Local, path)
}

// FileSet writes a file to the local disk.
func (m *localMedium) FileSet(path, content string) error {
	return io.Write(io.Local, path, content)
}

// Read reads a file from the local disk.
func (m *localMedium) Read(path string) (string, error) {
	return io.Read(io.Local, path)
}

// Write writes a file to the local disk.
func (m *localMedium) Write(path, content string) error {
	return io.Write(io.Local, path, content)
}

// EnsureDir creates a directory on the local disk.
func (m *localMedium) EnsureDir(path string) error {
	return io.EnsureDir(io.Local, path)
}

// IsFile checks if a path exists and is a file on the local disk.
func (m *localMedium) IsFile(path string) bool {
	return io.IsFile(io.Local, path)
}
