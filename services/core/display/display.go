package display

import (
	"embed"

	"github.com/letheanVPN/desktop/services/core/config"
	"github.com/letheanVPN/desktop/services/core/i18n"
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
// It is the main entry point for all display-related operations.
type Service struct {
	// --- Injected Dependencies ---
	app           *application.App
	configService *config.Service
	i18nService   *i18n.Service

	// --- Internal State ---
	brand         Brand
	assets        embed.FS
	windowHandles map[string]*application.WebviewWindow
}
