package main

import (
	"embed"
	_ "embed"
	"log"
	"time"

	"github.com/letheanVPN/desktop/services/blockchain"
	"github.com/letheanVPN/desktop/services/display"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	displayService := display.New(display.ClientHub, assets)

	app := application.New(application.Options{
		Name:        "desktop",
		Description: "Lethean Desktop",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Services: []application.Service{
			application.NewService(displayService),
			application.NewService(lthn.NewLetheanService()),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	// --- Create Main Window ---
	displayService.OpenWindow(app, "main", application.WebviewWindowOptions{
		Title: "Lethean Desktop",
		URL:   "#/editor/monaco", // Load the default Angular route
	})

	// Create and configure the app using the ConfigService AFTER the windows are created
	displayService.BuildMenu(app)

	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
