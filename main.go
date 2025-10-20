package main

import (
	"embed"
	"log"

	"github.com/letheanVPN/desktop/services/blockchain"
	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/letheanVPN/desktop/services/config"
	"github.com/letheanVPN/desktop/services/display"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// --- Initialize Core Services ---

	// 1. Create Config Service FIRST.
	configService, err := config.NewService()
	if err != nil {
		log.Fatalf("Fatal: Failed to initialize config service: %v", err)
	}
	cfg := configService.Get()

	// 2. Create other services, injecting dependencies into them.
	//cryptService := crypt.NewService(cfg) // Using the new standardized constructor
	//workspaceService := workspace.NewService(cfg)
	// Pass the embedded assets directly to the DisplayService.
	displayService := display.NewService(cfg, display.ClientHub, assets)

	letheanService := blockchain.NewService(cfg)

	// --- Initialize Wails Application ---
	app := application.New(application.Options{
		Name:        "Lethean Desktop",
		Description: "A private, decentralized, and secure desktop application",
		Services: []application.Service{
			application.NewService(displayService),
			application.NewService(letheanService),
			//application.NewService(cryptService),
			//application.NewService(workspaceService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})
	displayService.Setup(app)
	// --- Run Application ---
	displayService.OpenWindow(app, "main", application.WebviewWindowOptions{
		Title: "Lethean Code Editor",
		URL:   "#/editor/monaco", // Load the default Angular route
	})

	err = app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
