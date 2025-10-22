package display

import (
	"embed"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Brand defines the type for different application brands.
type Brand string

const (
	AdminHub     Brand = "admin-hub"
	ServerHub    Brand = "server-hub"
	GatewayHub   Brand = "gateway-hub"
	DeveloperHub Brand = "developer-hub"
	ClientHub    Brand = "client-hub"
)

// Service manages all OS-level UI interactions (menus, windows, tray).
// It is the main entry point for all display-related operations
// and is bound to the frontend.
type Service struct {
	brand         Brand
	assets        embed.FS
	windowHandles map[string]*application.WebviewWindow
	app           *application.App // Add this line
}
