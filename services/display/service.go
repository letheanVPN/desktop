package display

import (
	"embed"
	"fmt"

	"github.com/letheanVPN/desktop/services/config"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// NewService creates a new DisplayService.
func NewService(cfg *config.Config, brand Brand, assets embed.FS) *Service {
	return &Service{
		config:        cfg,
		brand:         brand,
		assets:        assets,
		windowHandles: make(map[string]*application.WebviewWindow),
	}
}

// Setup initializes the display service with the application instance.
func (s *Service) Setup(app *application.App) {
	s.app = app // Add this line
	s.buildMenu(app)
	s.setupTray(app)
}

func (s *Service) ShowEnvironmentDialog() {
	envInfo := s.app.Env.Info() // Use s.app

	details := fmt.Sprintf(`Environment Information:

Operating System: %s
Architecture: %s
Debug Mode: %t

Dark Mode: %t

Platform Information:`,
		envInfo.OS,
		envInfo.Arch,
		envInfo.Debug,
		s.app.Env.IsDarkMode()) // Use s.app

	// Add platform-specific details
	for key, value := range envInfo.PlatformInfo {
		details += fmt.Sprintf("\n%s: %v", key, value)
	}

	if envInfo.OSInfo != nil {
		details += fmt.Sprintf("\n\nOS Details:\nName: %s\nVersion: %s",
			envInfo.OSInfo.Name,
			envInfo.OSInfo.Version)
	}

	dialog := s.app.Dialog.Info() // Use s.app
	dialog.SetTitle("Environment Information")
	dialog.SetMessage(details)
	dialog.Show()
}
