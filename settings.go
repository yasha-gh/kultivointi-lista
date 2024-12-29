package main

import (
	"context"
	"encoding/json"
	"fmt"
	"kultivointi-lista/utils"
	"os"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Settings struct {
	ctx context.Context	`json:"-"`
	DbFilename string	`json:"dbFilename"` //Filename without extension
	DbFileDir string	`json:"dbFileDir"`
}


func (s *Settings) GetSelf() *Settings {
	return s
}

func (s *Settings) Save() error {
	log := utils.GetLogger()
	settingsBytes, err := json.Marshal(s)
	if err != nil {
		log.Error("Failed marshal settings", "err", err)
		return err
	}
	settingsFile, _ := GetSettingsFilePath()
	err = os.WriteFile(settingsFile, settingsBytes, 0644)
	if err != nil {
		log.Error("Failed to save settings to disk", "err", err)
		return err
	}
	return nil
}

func (s *Settings) Sync(settings Settings) error {
	runtime.EventsEmit(s.ctx, "settings_sync", s)
	if settings.DbFilename == "" || settings.DbFileDir == "" {
		return fmt.Errorf("Settings sync: Missing required data: DbFilename: %v, DbFileDir: %v", settings.DbFilename, settings.DbFileDir)
	}
	s.DbFilename = settings.DbFilename
	s.DbFileDir = settings.DbFileDir
	return nil
}

func (s *Settings) OnSync() error {
	log := utils.GetLogger()
	log.Info("settings OnSync listener started")
	runtime.EventsOn(s.ctx, "settings_sync", func(data ...interface{}) {
		utils.PrettyPrint(data)
		fmt.Println("settings sync event received")
	})
	return nil
}

func (s *Settings) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *Settings) GetContext() context.Context {
	return s.ctx
}

func (s *Settings) GetDbAbsolutePath() string {
	log := utils.Logger
	if s.DbFilename == "" {
		log.Fatal("GetDbAbsolutePath: missing data to get DB File absolute path", "DbFilePath", s.DbFileDir, "DbFilename", s.DbFilename)
	}
	return fmt.Sprintf("%s/%s.kvldb",s.DbFileDir, s.DbFilename)
}

func NewSettings(ctx context.Context) *Settings {
	log := utils.Logger
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal("NewSettings: Failed to get OS Cache dir")
	}

	dbFileDir := fmt.Sprintf("%s/kultivointi-lista", cacheDir)
	err = utils.CreateDirAll(dbFileDir)
	if err != nil {
		log.Fatal("Failed to create app cache dir", "err", err)
	}
	return &Settings{
		DbFileDir: dbFileDir,
		DbFilename: "kultivointi-lista",
		ctx: ctx,
	}
}

func GetSettingsFilePath() (filePath string, err error) {
	log := utils.Logger
	osCacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Error("GetSettingsFilePath: failed to get OS user cache dir", "err", err)
		return "", err
	}
	settingsFilePath := fmt.Sprintf("%s/kultivointi-lista/settings.json", osCacheDir)
	err = utils.PathExists(settingsFilePath)
	return settingsFilePath, err
}

func LoadSettings(ctx context.Context) (settings *Settings, err error) {
	log := utils.Logger
	settingsFile, err := GetSettingsFilePath()
	if err != nil {
		if os.IsNotExist(err) {
			settings = NewSettings(ctx)
			log.Info("LoadSettings: settings file not found, creating new")
			err = nil
			return
		} else {
			log.Error("LoadSettings: failed to load settings", "err", err)
			return
		}
	}
	fileBytes, err := os.ReadFile(settingsFile)
	if err != nil {
		log.Error("LoadSettings: Failed to read settings from disk","err", err)
		return
	}
	err = json.Unmarshal(fileBytes, &settings)
	if err != nil {
		log.Error("LoadSettings: Failed to unmarshal settings file", "err", err)
		return
	}
	log.Info("LoadSettings: Settings loaded from disk")
	return
}
