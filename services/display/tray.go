package display

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// setupTray configures and creates the system tray icon and menu.
func (s *Service) setupTray(app *application.App) {
	// Use the 'assets' field that was injected into the service.
	//iconBytes, _ := s.assets.ReadFile("frontend/dist/favicon.ico")

	systray := app.SystemTray.New()
	systray.SetTooltip("Lethean Desktop")
	systray.SetLabel("LTHN")
	//systray.SetIcon(iconBytes)

	// Create a hidden window for the system tray menu to interact with
	trayWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "System Tray Status",
		URL:       "/#/system-tray",
		Width:     400,
		Frameless: true,
		Hidden:    true,
	})
	systray.AttachWindow(trayWindow)

	// --- Build Tray Menu ---
	trayMenu := app.Menu.New()
	trayMenu.Add("Open").OnClick(func(ctx *application.Context) {
		for _, window := range app.Window.GetAll() {
			window.Show()
		}
	})
	trayMenu.Add("Close").OnClick(func(ctx *application.Context) {
		for _, window := range app.Window.GetAll() {
			window.Hide()
		}
	})

	// Add brand-specific menu items
	switch s.brand {
	case AdminHub:
		trayMenu.Add("Manage Workspace").OnClick(func(ctx *application.Context) { /* TODO */ })
	case ServerHub:
		trayMenu.Add("Server Control").OnClick(func(ctx *application.Context) { /* TODO */ })
	case GatewayHub:
		trayMenu.Add("Routing Table").OnClick(func(ctx *application.Context) { /* TODO */ })
	case DeveloperHub:
		trayMenu.Add("Debug Console").OnClick(func(ctx *application.Context) { /* TODO */ })
	case ClientHub:
		trayMenu.Add("Connect").OnClick(func(ctx *application.Context) { /* TODO */ })
		trayMenu.Add("Disconnect").OnClick(func(ctx *application.Context) { /* TODO */ })
	}

	trayMenu.AddSeparator()
	trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	systray.SetMenu(trayMenu)
}
