package docs

import (
	"embed"

	"github.com/letheanVPN/desktop/services/core/display"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// displayer is an interface that defines the functionality docs needs from a display service.
// This avoids a direct dependency on the display package or the core package.
type displayer interface {
	OpenWindow(name string, options application.WebviewWindowOptions) (*application.WebviewWindow, error)
}

// Service manages the documentation display and serving of assets.
type Service struct {
	// --- Injected Dependencies ---
	app            *application.App
	displayService *display.Service // Depends on the local interface, not a concrete type from another package.

	// --- Internal State ---
	assets embed.FS
}

//go:embed all:static/**/*
var docsStatic embed.FS
