package main

import (
	"PoCopilot/backend/common"
	"PoCopilot/backend/services/handler"
	"embed"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	initDir := common.GetConfigDir()
	logFile := filepath.Join(initDir, common.AppLogFile)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "PoCopilot",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: handler.NewFileLoader(),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Logger:             common.NewLogger(logFile),
		LogLevel:           logger.DEBUG,
		LogLevelProduction: logger.DEBUG,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
