package docs

import (
	"embed"

	"github.com/letheanVPN/desktop/services/core/display"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// Service manages the documentation display and serving of assets.
type Service struct {
	// --- Injected Dependencies ---
	app            *application.App
	displayService *display.Service

	// --- Internal State ---
	assets embed.FS
}

//go:embed all:static/**/*
var docsStatic embed.FS
