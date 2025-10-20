package display

import (
	"embed"

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

// Startup is called by Wails when the application starts.
// It orchestrates the setup of all UI elements.
func (s *Service) Setup(app *application.App) {
	s.buildMenu(app)
	s.setupTray(app)
}
