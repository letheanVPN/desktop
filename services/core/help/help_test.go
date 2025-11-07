package help

import (
	"testing"

	"github.com/Snider/Core/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// MockDisplay is a mock implementation of the core.Display interface.
type MockDisplay struct {
	ShowCalled bool
}

func (m *MockDisplay) Show() error {
	m.ShowCalled = true
	return nil
}

func (m *MockDisplay) ShowAt(anchor string) error {
	m.ShowCalled = true
	return nil
}

func (m *MockDisplay) Hide() error                                { return nil }
func (m *MockDisplay) HideAt(anchor string) error                 { return nil }
func (m *MockDisplay) OpenWindow(opts ...core.WindowOption) error { return nil }

// MockCore is a mock implementation of the *core.Core type.
type MockCore struct {
	Core         *core.Core
	ActionCalled bool
	ActionMsg    core.Message
}

// ACTION matches the signature required by RegisterAction.
func (m *MockCore) ACTION(c *core.Core, msg core.Message) error {
	m.ActionCalled = true
	m.ActionMsg = msg
	return nil
}

func setupService(t *testing.T) (*Service, *MockCore, *MockDisplay) {
	s, err := New()
	assert.NoError(t, err)

	app := application.New(application.Options{})
	c, err := core.New(core.WithWails(app))
	assert.NoError(t, err)
	mockCore := &MockCore{Core: c}
	mockDisplay := &MockDisplay{}

	s.Runtime = core.NewRuntime(c, Options{})
	s.display = mockDisplay
	// Register our mock handler. When the real s.Core().ACTION is called,
	// our mock handler will be executed.
	c.RegisterAction(mockCore.ACTION)

	return s, mockCore, mockDisplay
}

func TestNew(t *testing.T) {
	s, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, s)
}

func TestShow(t *testing.T) {
	s, mockCore, _ := setupService(t)

	err := s.Show()
	assert.NoError(t, err)
	assert.True(t, mockCore.ActionCalled)

	msg, ok := mockCore.ActionMsg.(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "display.open_window", msg["action"])
	assert.Equal(t, "help", msg["name"])
}

func TestShowAt(t *testing.T) {
	s, mockCore, _ := setupService(t)

	err := s.ShowAt("test-anchor")
	assert.NoError(t, err)
	assert.True(t, mockCore.ActionCalled)

	msg, ok := mockCore.ActionMsg.(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "display.open_window", msg["action"])
	assert.Equal(t, "help", msg["name"])

	opts, ok := msg["options"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "/#test-anchor", opts["URL"])
}

func TestHandleIPCEvents_ServiceStartup(t *testing.T) {
	s, _, _ := setupService(t)
	err := s.HandleIPCEvents(s.Core(), core.ActionServiceStartup{})
	assert.NoError(t, err)
}
