package display

import (
	"context"
	"embed"
	"fmt"

	"github.com/letheanVPN/desktop/services/core/config"
	"github.com/letheanVPN/desktop/services/core/i18n"
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

// Setup initializes the display service with the application instance and other core services.
func (s *Service) Setup(app *application.App, configService *config.Service, i18nService *i18n.Service) {
	s.app = app
	s.configService = configService
	s.i18nService = i18nService

	s.analyzeScreens()
	s.monitorScreenChanges()
	s.buildMenu()
	s.setupTray()
}

func (s *Service) analyzeScreens() {
	s.app.Logger.Info("Screen analysis", "count", len(s.app.Screen.GetAll()))

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

	for i, screen := range s.app.Screen.GetAll() {
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
		s.app.Logger.Info("Updated screen count", "count", len(s.app.Screen.GetAll()))

		// Could reposition windows here if needed
	})
}

func (s *Service) ShowEnvironmentDialog() {
	envInfo := s.app.Env.Info()

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

	dialog := s.app.Dialog.Info()
	dialog.SetTitle("Environment Information")
	dialog.SetMessage(details)
	dialog.Show()
}

func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	// Check the IsNew flag from the config service, which is the single source of truth.
	if s.configService.Get().IsNew {
		// If the config was just created, open the setup window.
		s.OpenWindow("main", application.WebviewWindowOptions{
			Title: s.i18nService.Translate("app.setup.title"),
			URL:   "#/setup",
		})
	} else {
		// If the config already existed, open the main window with the default route.
		defaultRoute, err := s.configService.Get().Key("DefaultRoute")
		if err != nil {
			defaultRoute = "/" // Fallback to a safe default if the key is somehow missing.
			s.app.Logger.Error("Could not get DefaultRoute from config, using fallback: %v", err)
		}

		s.OpenWindow("main", application.WebviewWindowOptions{
			Title:  s.i18nService.Translate("app.title"),
			Height: 900,
			Width:  1280,
			URL:    "#" + defaultRoute.(string),
		})
	}
	return nil
}
