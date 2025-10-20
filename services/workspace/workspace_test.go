package workspace

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/letheanVPN/desktop/services/config"
	"github.com/stretchr/testify/assert"
)

// MockMedium implements the Medium interface for testing purposes.
type MockMedium struct {
	Files map[string]string
	Dirs  map[string]bool
}

func NewMockMedium() *MockMedium {
	return &MockMedium{
		Files: make(map[string]string),
		Dirs:  make(map[string]bool),
	}
}

func (m *MockMedium) FileGet(path string) (string, error) {
	content, ok := m.Files[path]
	if !ok {
		return "", assert.AnError // Simulate file not found error
	}
	return content, nil
}

func (m *MockMedium) FileSet(path, content string) error {
	m.Files[path] = content
	return nil
}

func (m *MockMedium) EnsureDir(path string) error {
	m.Dirs[path] = true
	return nil
}

func (m *MockMedium) IsFile(path string) bool {
	_, ok := m.Files[path]
	return ok
}

func TestNewService(t *testing.T) {
	mockConfig := &config.Config{} // You might want to mock this further if its behavior is critical
	mockMedium := NewMockMedium()

	service := NewService(mockConfig, mockMedium)

	assert.NotNil(t, service)
	assert.Equal(t, mockConfig, service.config)
	assert.Equal(t, mockMedium, service.medium)
	assert.NotNil(t, service.workspaceList)
	assert.Nil(t, service.activeWorkspace) // Initially no active workspace
}

func TestServiceStartup(t *testing.T) {
	mockConfig := &config.Config{
		WorkspacesDir: "/tmp/workspaces",
	}

	// Test case 1: list.json exists and is valid
	t.Run("existing valid list.json", func(t *testing.T) {
		mockMedium := NewMockMedium()

		// Prepare a mock workspace list
		expectedWorkspaceList := map[string]string{
			"workspace1": "pubkey1",
			"workspace2": "pubkey2",
		}
		listContent, _ := json.MarshalIndent(expectedWorkspaceList, "", "  ")

		listPath := filepath.Join(mockConfig.WorkspacesDir, listFile)
		mockMedium.FileSet(listPath, string(listContent))

		service := NewService(mockConfig, mockMedium)
		err := service.ServiceStartup()

		assert.NoError(t, err)
		assert.Equal(t, expectedWorkspaceList, service.workspaceList)
		assert.NotNil(t, service.activeWorkspace)
		assert.Equal(t, defaultWorkspace, service.activeWorkspace.Name)
		assert.Equal(t, filepath.Join(mockConfig.WorkspacesDir, defaultWorkspace), service.activeWorkspace.Path)
	})

	// Test case 2: list.json does not exist
	t.Run("no list.json", func(t *testing.T) {
		mockMedium := NewMockMedium() // Fresh medium with no files

		service := NewService(mockConfig, mockMedium)
		err := service.ServiceStartup()

		assert.NoError(t, err)
		assert.NotNil(t, service.workspaceList)
		assert.Empty(t, service.workspaceList) // Should be empty if no list.json
		assert.NotNil(t, service.activeWorkspace)
		assert.Equal(t, defaultWorkspace, service.activeWorkspace.Name)
		assert.Equal(t, filepath.Join(mockConfig.WorkspacesDir, defaultWorkspace), service.activeWorkspace.Path)
	})

	// Test case 3: list.json exists but is invalid
	t.Run("invalid list.json", func(t *testing.T) {
		mockMedium := NewMockMedium()

		listPath := filepath.Join(mockConfig.WorkspacesDir, listFile)
		mockMedium.FileSet(listPath, "{invalid json") // Invalid JSON

		service := NewService(mockConfig, mockMedium)
		err := service.ServiceStartup()

		assert.NoError(t, err) // Error is logged, but startup continues
		assert.NotNil(t, service.workspaceList)
		assert.Empty(t, service.workspaceList) // Should be empty if invalid list.json
		assert.NotNil(t, service.activeWorkspace)
		assert.Equal(t, defaultWorkspace, service.activeWorkspace.Name)
		assert.Equal(t, filepath.Join(mockConfig.WorkspacesDir, defaultWorkspace), service.activeWorkspace.Path)
	})
}
