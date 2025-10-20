package workspace

import (
	"github.com/letheanVPN/desktop/services/config"
)

const (
	defaultWorkspace = "default"
	listFile         = "list.json"
)

// Workspace represents a user's workspace.
type Workspace struct {
	Name string
	Path string
}

// Service manages user workspaces.
type Service struct {
	config          *config.Config
	activeWorkspace *Workspace
	workspaceList   map[string]string // Maps Workspace ID to Public Key
	medium          Medium
}

// Medium defines the interface for a workspace storage medium.
type Medium interface {
	FileGet(path string) (string, error)
	FileSet(path, content string) error
	EnsureDir(path string) error
	IsFile(path string) bool
}
