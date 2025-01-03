package main

import (
	"context"
	"embed"
	"log"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"kultivointi-lista/db"
	"kultivointi-lista/list"
	"kultivointi-lista/utils"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	utils.Logger = utils.NewLogger()
	if utils.Logger == nil {
		log.Fatal("Failed create new logger")
	}
	// Create an instance of the app structure
	app := NewApp()

	settings, err := LoadSettings(context.Background())
	if err != nil {
		log := utils.GetLogger()
		log.Fatal("Failed to load settings", "err", err)
	}

	parser := &list.ListParser{}
	// list := &list.List{
	// 	MainList: make([]*list.ListItem, 0),
	// 	Filtered: make([]*list.ListItem, 0),
	// }
	list := list.NewList(nil)
	parser.List = list
	// Create application with options
	err = wails.Run(&options.App{
		Title:  "kultivointi-lista",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		// OnStartup:        app.startup,
		OnStartup: func (ctx context.Context) {
			log := utils.GetLogger()
			envInfo := runtime.Environment(ctx)
			if envInfo.BuildType == "production" {
				isDev := false
				utils.IsDevMode = &isDev
			} else {
				isDev := true
				utils.IsDevMode = &isDev
			}
			settings.SetContext(ctx)
			dbFile := settings.GetDbAbsolutePath()
			dbPool, err := db.TursoConnect(dbFile)
			if err != nil {
				log.Fatal("OnStartup: failed to init db pool", "err", err)
			}
			ctx = context.WithValue(ctx,"dbPool", dbPool)
			settings.SetContext(ctx)
			go settings.OnSync()
			// settings.Sync(*settings)
			db.CreateTables(dbPool)
			parser.SetContext(ctx)
			list.SetContext(ctx)
			go list.ListChannels()
			app.ctx = ctx
		},
		Bind: []interface{}{
			app,
			settings,
			parser,
			list,
		},
		OnDomReady: func(ctx context.Context) {
			log := utils.GetLogger()
			log.Info("Wails OnDomReady")
			// settings.Sync(*settings)
		},
	})

	if err != nil {
		log := utils.GetLogger()
		log.Error("Wails Run error", "err", err)
		// println("Error:", err.Error())
	}
}
