package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/letheanVPN/desktop/services/blockchain"
	"github.com/letheanVPN/desktop/services/config"
	"github.com/letheanVPN/desktop/services/display"
	"github.com/letheanVPN/desktop/services/docs"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
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

	docsService, err := docs.NewService(displayService)
	if err != nil {
		log.Fatalf("Fatal: Failed to initialize docs service: %v", err)
	}

	// --- Initialize Wails Application ---
	app := application.New(application.Options{
		Name:        "Lethean Desktop",
		Description: "A private, decentralized, and secure desktop application",
		Services: []application.Service{
			application.NewService(displayService),
			application.NewService(letheanService),
			application.NewService(configService),
			application.NewServiceWithOptions(docsService, application.ServiceOptions{
				Route: "docs",
			}),
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

	// OS specific application events

	// Platform agnostic events
	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
		app.Logger.Info("Application started!")
	})

	app.Event.OnApplicationEvent(events.Windows.SystemThemeChanged, func(event *application.ApplicationEvent) {
		app.Logger.Info("System theme changed!")
		if event.Context().IsDarkMode() {
			app.Logger.Info("System is now using dark mode!")
		} else {
			app.Logger.Info("System is now using light mode!")
		}
	})
	displayService.Setup(app)
	// --- Run Application ---
	configFilePath := filepath.Join(cfg.ConfigDir, "config.json")
	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		displayService.OpenWindow("main", application.WebviewWindowOptions{
			Title: "Desktop Setup",
			URL:   "#/setup",
		})
	} else {
		//displayService.OpenWindow("main", application.WebviewWindowOptions{
		//	Title:  "Desktopfdfd",
		//	Height: 900,
		//	Width:  1280,
		//	URL:    "/docs/",
		//})
		displayService.OpenWindow("main", application.WebviewWindowOptions{
			Title:  "Desktop",
			Height: 900,
			Width:  1280,
			URL:    "#" + cfg.DefaultRoute,
		})
	}

	err = app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
