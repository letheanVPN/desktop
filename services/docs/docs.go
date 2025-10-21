package docs

import (
	"embed"

	"github.com/letheanVPN/desktop/services/display"
)

// Docs holds the resolved paths for the application.
type Config struct {
	DataDir string // For user-specific data files.
}
type Service struct {
	config  *Config
	display *display.Service
}

//go:embed all:static/**/*
var docsStatic embed.FS
