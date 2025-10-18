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

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Window 1",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	// Create and configure the app using the ConfigService AFTER the window is created
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
