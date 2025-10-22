package display

import (
	"embed"
	"fmt"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// NewService creates a new DisplayService.
func NewService(brand Brand, assets embed.FS) *Service {
	return &Service{
		brand:         brand,
		assets:        assets,
		windowHandles: make(map[string]*application.WebviewWindow),
	}
}

// Setup initializes the display service with the application instance.
func (s *Service) Setup(app *application.App) {
	s.app = app
	s.analyzeScreens()
	s.monitorScreenChanges()
	s.buildMenu()
	s.setupTray()
}

func (s *Service) analyzeScreens() {
	screens := s.app.Screen.GetAll()
	s.app.Logger.Info("Screen analysis", "count", len(screens))

	primary := s.app.Screen.GetPrimary()
	if primary != nil {
		s.app.Logger.Info("Primary screen",
			"name", primary.Name,
			"size", fmt.Sprintf("%dx%d", primary.Size.Width, primary.Size.Height),
			"scaleFactor", primary.ScaleFactor,
			"workArea", primary.WorkArea,
		)
		scaleFactor := primary.ScaleFactor

		switch {
		case scaleFactor == 1.0:
			s.app.Logger.Info("Standard DPI display", "screen", primary.Name)
		case scaleFactor == 1.25:
			s.app.Logger.Info("125% scaled display", "screen", primary.Name)
		case scaleFactor == 1.5:
			s.app.Logger.Info("150% scaled display", "screen", primary.Name)
		case scaleFactor == 2.0:
			s.app.Logger.Info("High DPI display (200%)", "screen", primary.Name)
		default:
			s.app.Logger.Info("Custom scale display",
				"screen", primary.Name,
				"scale", scaleFactor,
			)
		}
	} else {
		s.app.Logger.Info("No primary screen found")
	}

	for i, screen := range screens {
		s.app.Logger.Info("Screen details",
			"index", i,
			"name", screen.Name,
			"primary", screen.IsPrimary,
			"bounds", screen.Bounds,
			"scaleFactor", screen.ScaleFactor,
		)
	}
}

func (s *Service) monitorScreenChanges() {
	// Monitor for screen configuration changes
	s.app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(event *application.ApplicationEvent) {
		s.app.Logger.Info("Screen configuration changed")

		// Re-analyze screens
		screens := s.app.Screen.GetAll()
		s.app.Logger.Info("Updated screen count", "count", len(screens))

		// Could reposition windows here if needed
	})
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
