package main

import (
	"embed"
	"log"

	"github.com/letheanVPN/desktop/services/blockchain"
	"github.com/letheanVPN/desktop/services/core"
	"github.com/letheanVPN/desktop/services/mining"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// --- Phase 1: Initialize Core Service Singleton ---
	// The core package is now responsible for instantiating all its services.
	// We pass the assets FS, as it's a compile-time dependency for some services.
	core.New(assets)

	// --- Initialize Wails Application ---
	app := application.New(application.Options{
		Name:        "Lethean Desktop",
		Description: "A private, decentralized, and secure desktop application",
		Services: []application.Service{
			// Register all core services with Wails via the clean core API.
			application.NewService(core.Config()),
			application.NewService(core.I18n()),
			application.NewService(core.Display()),
			application.NewService(core.Docs()),

			// Register other application services.
			application.NewService(blockchain.NewService()),
			application.NewService(mining.New()),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	// --- Phase 2: Wire Core Services ---
	// Pass the application instance to the core services that need it.
	core.Setup(app)

	// --- Run Application ---
	// The display service's ServiceStartup hook will handle opening the initial window.
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
