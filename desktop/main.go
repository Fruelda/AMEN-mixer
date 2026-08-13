package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// ============================================================
// FRONTEND ASSETS
// ============================================================

//go:embed all:frontend/dist
var assets embed.FS

// ============================================================
// MAIN
// ============================================================

func main() {
	app := NewApp()

	// ========================================================
	// WAILS APPLICATION
	// ========================================================

	err := wails.Run(&options.App{
		Title:  "desktop",
		Width:  1024,
		Height: 768,

		AssetServer: &assetserver.Options{
			Assets: assets,
		},

		BackgroundColour: &options.RGBA{
			R: 27,
			G: 38,
			B: 54,
			A: 1,
		},

		OnStartup:  app.startup,
		OnShutdown: app.shutdown,

		Bind: []interface{}{
			app,
		},
	})

	// ========================================================
	// ERROR
	// ========================================================

	if err != nil {
		println(
			"Error:",
			err.Error(),
		)
	}
}
