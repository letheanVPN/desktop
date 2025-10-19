package config

import (
	"embed"
	"fmt"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Brand defines the type for different application brands
type Brand string

const (
	AdminHub     Brand = "admin-hub"
	ServerHub    Brand = "server-hub"
	GatewayHub   Brand = "gateway-hub"
	DeveloperHub Brand = "developer-hub"
	ClientHub    Brand = "client-hub"
)

type ConfigService struct {
}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

// ConfigureApp applies settings to the main application object
func (s *ConfigService) ConfigureApp(app *application.App, brand Brand, assets embed.FS) {
	fmt.Printf("Configuring app for brand: %s\n", brand)
	// --- Setup System Tray ---
	iconBytes, _ := assets.ReadFile("frontend/browser/favicon.ico")
	systray := app.SystemTray.New()
	systray.SetTooltip("My Application Tooltip") // Windows
	systray.SetLabel("LTHN")                     // Windows
	systray.SetIcon(iconBytes)

	// Create a window
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "System Tray Status",
		URL:       "/#/system-tray", // Load the system tray Angular route
		Width:     400,
		Frameless: true,
		Hidden:    true,
	})
	systray.AttachWindow(window)
	appMenu := app.Menu.New()
	if runtime.GOOS == "darwin" {
		appMenu.AddRole(application.AppMenu)
	}
	appMenu.AddRole(application.FileMenu)
	appMenu.AddRole(application.ViewMenu)
	appMenu.AddRole(application.EditMenu)
	trayMenu := app.Menu.New()
	trayMenu.Add("Open").OnClick(func(ctx *application.Context) {
		windows := app.Window.GetAll()
		for _, window := range windows {
			window.Show()
		}
	})
	trayMenu.Add("Close").OnClick(func(ctx *application.Context) {
		windows := app.Window.GetAll()
		for _, window := range windows {
			window.Hide()
		}
	})
	workspace := appMenu.AddSubmenu("Workspace")
	workspace.Add("New").OnClick(func(ctx *application.Context) {})
	workspace.Add("List").OnClick(func(ctx *application.Context) {})
	// just sudo code demo to remind later me;
	switch brand {
	case AdminHub:
		trayMenu.Add("Manage Workspace").OnClick(func(ctx *application.Context) {})
	case ServerHub:
		trayMenu.Add("Server Control").OnClick(func(ctx *application.Context) {})
	case GatewayHub:
		trayMenu.Add("Routing Table").OnClick(func(ctx *application.Context) {})
	case DeveloperHub:
		appMenu.AddSubmenu("Developer")
		trayMenu.Add("Debug Console").OnClick(func(ctx *application.Context) {})
	case ClientHub:
		trayMenu.Add("Connect").OnClick(func(ctx *application.Context) {})
		trayMenu.Add("Disconnect").OnClick(func(ctx *application.Context) {})
	}

	trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	appMenu.AddRole(application.WindowMenu)
	appMenu.AddRole(application.HelpMenu)

	systray.SetMenu(trayMenu)
	app.Menu.Set(appMenu)

}
