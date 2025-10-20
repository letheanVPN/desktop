package workspace

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/letheanVPN/desktop/services/crypt"
	"github.com/letheanVPN/desktop/services/crypt/lib/openpgp"
	"github.com/letheanVPN/desktop/services/filesystem"
)

const (
	workspacesDir    = "workspaces"
	defaultWorkspace = "default"
	listFile         = "list.json"
)

// Workspace represents a user's workspace, containing their configuration and data.
type Workspace struct {
	Name string
	Path string
}

// WorkspaceService manages the active workspace.
type WorkspaceService struct {
	activeWorkspace *Workspace
	workspaceList   map[string]string // Maps Workspace ID to Public Key
}

// NewWorkspaceService creates a new WorkspaceService.
func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{
		workspaceList: make(map[string]string),
	}
}

// Startup initializes the default workspace and loads the workspace list.
func (s *WorkspaceService) Startup() error {
	// Ensure root workspaces directory exists
	if err := filesystem.EnsureDir(workspacesDir); err != nil {
		return fmt.Errorf("failed to ensure root workspaces directory: %w", err)
	}

	// Load the workspace list
	listPath := filepath.Join(workspacesDir, listFile)
	if filesystem.IsFile(listPath) {
		content, err := filesystem.Read(listPath)
		if err != nil {
			return fmt.Errorf("failed to read workspace list: %w", err)
		}
		if err := json.Unmarshal([]byte(content), &s.workspaceList); err != nil {
			return fmt.Errorf("failed to parse workspace list: %w", err)
		}
	}

	// For now, we still switch to a default workspace. This could be changed later.
	return s.SwitchWorkspace(defaultWorkspace)
}

// CreateWorkspace creates a new, obfuscated workspace for a user.
func (s *WorkspaceService) CreateWorkspace(identifier, password string) (string, error) {
	// 1. Generate obfuscated ID
	realName := crypt.Hash("lthn", identifier)
	workspaceID := crypt.Hash("lthn", fmt.Sprintf("workspace/%s", realName))
	workspacePath := filepath.Join(workspacesDir, workspaceID)

	// Check if workspace already exists to prevent overwriting
	if _, exists := s.workspaceList[workspaceID]; exists {
		return "", fmt.Errorf("workspace for this identifier already exists")
	}

	// 2. Create directory structure
	dirsToCreate := []string{"config", "log", "data", "files", "keys"}
	for _, dir := range dirsToCreate {
		if err := filesystem.EnsureDir(filepath.Join(workspacePath, dir)); err != nil {
			return "", fmt.Errorf("failed to create workspace directory '%s': %w", dir, err)
		}
	}

	// 3. Generate and store OpenPGP key pair
	keyPair, err := openpgp.CreateKeyPair(workspaceID, password)
	if err != nil {
		return "", fmt.Errorf("failed to create workspace key pair: %w", err)
	}

	keyFiles := map[string]string{
		filepath.Join(workspacePath, "keys", "key.pub"):  keyPair.PublicKey,
		filepath.Join(workspacePath, "keys", "key.priv"): keyPair.PrivateKey,
	}
	for path, content := range keyFiles {
		if err := filesystem.Write(path, content); err != nil {
			return "", fmt.Errorf("failed to write key file %s: %w", path, err)
		}
	}

	// 4. Update and save the central workspace list
	s.workspaceList[workspaceID] = keyPair.PublicKey
	listData, err := json.MarshalIndent(s.workspaceList, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal workspace list: %w", err)
	}

	listPath := filepath.Join(workspacesDir, listFile)
	if err := filesystem.Write(listPath, string(listData)); err != nil {
		return "", fmt.Errorf("failed to write workspace list file: %w", err)
	}

	return workspaceID, nil
}

// SwitchWorkspace changes the active workspace.
func (s *WorkspaceService) SwitchWorkspace(name string) error {
	// We need to check if the workspace exists in our list, except for the "default" one
	if name != defaultWorkspace {
		if _, exists := s.workspaceList[name]; !exists {
			return fmt.Errorf("workspace '%s' does not exist", name)
		}
	}

	path := filepath.Join(workspacesDir, name)
	if err := filesystem.EnsureDir(path); err != nil {
		return fmt.Errorf("failed to ensure workspace directory exists: %w", err)
	}

	s.activeWorkspace = &Workspace{
		Name: name,
		Path: path,
	}

	return nil
}

// GetActiveWorkspace returns the name of the currently active workspace.
func (s *WorkspaceService) GetActiveWorkspace() string {
	if s.activeWorkspace == nil {
		return ""
	}
	return s.activeWorkspace.Name
}

// WorkspaceFileGet retrieves a file from the active workspace.
func (s *WorkspaceService) WorkspaceFileGet(filename string) (string, error) {
	if s.activeWorkspace == nil {
		return "", fmt.Errorf("no active workspace")
	}

	path := filepath.Join(s.activeWorkspace.Path, filename)
	return filesystem.Read(path)
}

// WorkspaceFileSet writes a file to the active workspace.
func (s *WorkspaceService) WorkspaceFileSet(filename, content string) error {
	if s.activeWorkspace == nil {
		return fmt.Errorf("no active workspace")
	}

	path := filepath.Join(s.activeWorkspace.Path, filename)
	return filesystem.Write(path, content)
}
