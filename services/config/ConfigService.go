package config

import (
	"embed"
	"fmt"

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
	iconBytes, _ := assets.ReadFile("frontend/dist/wails.png")
	systray := app.SystemTray.New()
	systray.SetTooltip("My Application Tooltip") // Windows
	systray.SetIcon(iconBytes)

	// Create a window
	//window := app.Window.NewWithOptions(application.WebviewWindowOptions{
	//	Title:     "Window 1",
	//	URL:       "/",
	//	Frameless: true,
	//})
	// systray.AttachWindow(window)

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

	// just sudo code demo to remind later me;
	switch brand {
	case AdminHub:
		trayMenu.Add("Manage Workspace").OnClick(func(ctx *application.Context) {})
	case ServerHub:
		trayMenu.Add("Server Control").OnClick(func(ctx *application.Context) {})
	case GatewayHub:
		trayMenu.Add("Routing Table").OnClick(func(ctx *application.Context) {})
	case DeveloperHub:
		trayMenu.Add("Debug Console").OnClick(func(ctx *application.Context) {})
	case ClientHub:
		trayMenu.Add("Connect").OnClick(func(ctx *application.Context) {})
		trayMenu.Add("Disconnect").OnClick(func(ctx *application.Context) {})
	}

	trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})
	systray.SetMenu(trayMenu)

}
