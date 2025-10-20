package workspace

import "github.com/letheanVPN/desktop/services/filesystem"

// localMedium implements the Medium interface for the local disk.
type localMedium struct{}

// newLocalMedium creates a new instance of the local storage medium.
func newLocalMedium() Medium {
	return &localMedium{}
}

// FileGet reads a file from the local disk.
func (m *localMedium) FileGet(path string) (string, error) {
	return filesystem.Read(filesystem.Local, path)
}

// FileSet writes a file to the local disk.
func (m *localMedium) FileSet(path, content string) error {
	return filesystem.Write(filesystem.Local, path, content)
}

// EnsureDir creates a directory on the local disk.
func (m *localMedium) EnsureDir(path string) error {
	return filesystem.EnsureDir(filesystem.Local, path)
}

// IsFile checks if a path exists and is a file on the local disk.
func (m *localMedium) IsFile(path string) bool {
	return filesystem.IsFile(filesystem.Local, path)
}
