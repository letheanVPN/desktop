package docs

import (
	"embed"

	"github.com/letheanVPN/desktop/services/core/display"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// Config holds the resolved paths for the application.
type Config struct {
	DataDir string // For user-specific data files.
}

// Service manages the documentation display and serving of assets.
type Service struct {
	// --- Injected Dependencies ---
	app            *application.App
	displayService display.Service

	// --- Internal State ---
	config *Config
}

//go:embed all:static/**/*
var docsStatic embed.FS
