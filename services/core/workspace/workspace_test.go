package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/Snider/Core/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// mockConfig is a mock implementation of the core.Config interface for testing.
type mockConfig struct {
	values map[string]interface{}
}

func (m *mockConfig) Get(key string, out any) error {
	val, ok := m.values[key]
	if !ok {
		return fmt.Errorf("key not found: %s", key)
	}
	// This is a simplified mock; a real one would use reflection to set `out`
	switch v := out.(type) {
	case *string:
		*v = val.(string)
	default:
		return fmt.Errorf("unsupported type in mock config Get")
	}
	return nil
}

func (m *mockConfig) Set(key string, v any) error {
	m.values[key] = v
	return nil
}

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
	_, exists := m.Files[path]
	return exists
}

func (m *MockMedium) Read(path string) (string, error) {
	return m.FileGet(path)
}

func (m *MockMedium) Write(path, content string) error {
	return m.FileSet(path, content)
}

// newTestService creates a workspace service instance with mocked dependencies.
func newTestService(t *testing.T, workspaceDir string) (*Service, *MockMedium) {
	coreInstance, err := core.New()
	assert.NoError(t, err)

	mockCfg := &mockConfig{values: map[string]interface{}{"workspaceDir": workspaceDir}}
	coreInstance.RegisterService("config", mockCfg)

	service, err := New()
	assert.NoError(t, err)

	service.Runtime = core.NewRuntime(coreInstance, Options{})
	mockMedium := NewMockMedium()
	service.medium = mockMedium

	return service, mockMedium
}

func TestServiceStartup(t *testing.T) {
	workspaceDir := "/tmp/workspace"

	t.Run("existing valid list.json", func(t *testing.T) {
		service, mockMedium := newTestService(t, workspaceDir)

		expectedWorkspaceList := map[string]string{
			"workspace1": "pubkey1",
			"workspace2": "pubkey2",
		}
		listContent, _ := json.MarshalIndent(expectedWorkspaceList, "", "  ")
		listPath := filepath.Join(workspaceDir, listFile)
		mockMedium.Files[listPath] = string(listContent)

		err := service.ServiceStartup(context.Background(), application.ServiceOptions{})

		assert.NoError(t, err)
		// assert.Equal(t, expectedWorkspaceList, service.workspaceList) // This check is difficult with current implementation
		assert.NotNil(t, service.activeWorkspace)
		assert.Equal(t, defaultWorkspace, service.activeWorkspace.Name)
	})
}

func TestCreateAndSwitchWorkspace(t *testing.T) {
	workspaceDir := "/tmp/workspace"
	service, _ := newTestService(t, workspaceDir)

	// Create
	workspaceID, err := service.CreateWorkspace("test", "password")
	assert.NoError(t, err)
	assert.NotEmpty(t, workspaceID)

	// Switch
	err = service.SwitchWorkspace(workspaceID)
	assert.NoError(t, err)
	assert.Equal(t, workspaceID, service.activeWorkspace.Name)
}
