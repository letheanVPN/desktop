package main

import (
	"embed"
	_ "embed"
	"log"
	"time"

	"github.com/letheanVPN/desktop/services/config"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	app := application.New(application.Options{
		Name:        "desktop",
		Description: "Lethean Desktop",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	// --- Create Main Window ---
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Lethean Desktop",
		URL:   "/", // Load the default Angular route
	})

	// --- Create Blockchain Status Window ---
	//app.Window.NewWithOptions(application.WebviewWindowOptions{
	//	Title: "Blockchain Status",
	//	URL:   "/#/blockchain", // Load the blockchain Angular route
	//	Width: 400,
	//	Height: 300,
	//})

	// Create and configure the app using the ConfigService AFTER the windows are created
	configService := config.NewConfigService()
	configService.ConfigureApp(app, config.ClientHub, assets)

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
