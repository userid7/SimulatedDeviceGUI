package main

import (
	"context"
	"embed"

	pm "SimulatedDeviceGUI/device/pm"
	hfreader "SimulatedDeviceGUI/device/hfreader"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	// app := NewApp()

	hfapp := hfreader.NewReaderApp()
	pmapp := pm.NewPMApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "SimulatedDeviceGUI",
		Width:            852,
		Height:           768,
		Assets:           assets,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			db := initDB()
			hfapp.Startup(ctx, db)
			pmapp.Startup(ctx, db)
		},
		Bind: []interface{}{
			hfapp,
			pmapp,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("deviceManager.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
