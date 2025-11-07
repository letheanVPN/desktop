package help

import (
	"context"
	"embed"
	"fmt"

	"github.com/Snider/Core/pkg/core"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:public/*
var helpStatic embed.FS

// Options holds configuration for the help service.
type Options struct{}

// Service manages the in-app help system.
type Service struct {
	*core.Runtime[Options]
	config  core.Config
	display core.Display
	assets  embed.FS
}

// newHelpService contains the common logic for initialising a Service struct.
func newHelpService() (*Service, error) {
	return &Service{
		assets: helpStatic,
	}, nil
}

// New is the constructor for static dependency injection.
// It creates a Service instance without initialising the core.Runtime field.
// Dependencies are passed directly here.
func New() (*Service, error) {
	s, err := newHelpService()
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Register is the constructor for dynamic dependency injection (used with core.WithService).
// It creates a Service instance and initialises its core.Runtime field.
// Dependencies are injected during ServiceStartup.
func Register(c *core.Core) (any, error) {
	s, err := newHelpService()
	if err != nil {
		return nil, err
	}
	s.Runtime = core.NewRuntime(c, Options{})
	return s, nil
}

// HandleIPCEvents processes IPC messages, including injecting dependencies on startup.
func (s *Service) HandleIPCEvents(c *core.Core, msg core.Message) error {
	switch m := msg.(type) {
	case core.ActionServiceStartup:
		return s.ServiceStartup(context.Background(), application.ServiceOptions{})
	default:
		c.App.Logger.Error("Help: Unknown message type", "type", fmt.Sprintf("%T", m))
	}
	return nil
}

// ServiceStartup is called when the app starts, after dependencies are injected.
func (s *Service) ServiceStartup(context.Context, application.ServiceOptions) error {
	s.Core().App.Logger.Info("Help service started")
	return nil
}

// Show displays the help window.
func (s *Service) Show() error {
	if s.display == nil {
		return fmt.Errorf("display service not initialized")
	}
	if s.Core() == nil {
		return fmt.Errorf("core runtime not initialized")
	}
	msg := map[string]any{
		"action": "display.open_window",
		"name":   "help",
		"options": map[string]any{
			"Title":  "Help",
			"Width":  800,
			"Height": 600,
		},
	}

	return s.Core().ACTION(msg)
}

// ShowAt displays a specific section of the help documentation.
func (s *Service) ShowAt(anchor string) error {
	if s.display == nil {
		return fmt.Errorf("display service not initialized")
	}
	if s.Core() == nil {
		return fmt.Errorf("core runtime not initialized")
	}
	msg := map[string]any{
		"action": "display.open_window",
		"name":   "help",
		"options": map[string]any{
			"Title":  "Help",
			"Width":  800,
			"Height": 600,
			"URL":    fmt.Sprintf("/#%s", anchor),
		},
	}
	return s.Core().ACTION(msg)
}

// Ensure Service implements the core.Help interface.
var _ core.Help = (*Service)(nil)
