package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/letheanVPN/desktop/services/blockchain"
	"github.com/letheanVPN/desktop/services/core"
	"github.com/letheanVPN/desktop/services/display"
	"github.com/letheanVPN/desktop/services/docs"
	"github.com/letheanVPN/desktop/services/mining"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// --- Initialize Core Services ---

	// Initialize the CoreService singleton. This must be called once at startup.
	core.New()
	cfg := core.Config().Get()

	// Create other services, injecting dependencies into them.
	displayService := display.NewService(display.ClientHub, assets)
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
			application.NewService(core.Config()),
			application.NewService(core.I18n()),
			application.NewService(displayService),
			application.NewService(letheanService),
			application.NewServiceWithOptions(docsService, application.ServiceOptions{
				Route: "docs",
			}),
			application.NewService(mining.New()),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	// Platform agnostic events
	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
		app.Logger.Info(core.I18n().Translate("app.boot.loaded-runtime"))
	})

	displayService.Setup(app)

	// --- Log startup message ---
	app.Logger.Info(core.I18n().Translate("app.boot.start-runtime"))

	// --- Run Application ---
	configFilePath := filepath.Join(cfg.ConfigDir, "config.json")
	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		displayService.OpenWindow("main", application.WebviewWindowOptions{
			Title: core.I18n().Translate("app.setup.title"),
			URL:   "#/setup",
		})
	} else {
		displayService.OpenWindow("main", application.WebviewWindowOptions{
			Title:  core.I18n().Translate("app.title"),
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
