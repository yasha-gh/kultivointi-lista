package db

import (
	"database/sql"
	_ "embed"
	"kultivointi-lista/utils"
	"fmt"
)

//go:embed sql/list_item_schema.sql
var listItem string

//go:embed sql/list_item_titles_schema.sql
var listItemTitles string

//go:embed sql/media_schema.sql
var media string

// //go:embed sql/season_schema.sql
// var season string

//go:embed sql/sites_schema.sql
var sites string


//go:embed sql/set_pragma.sql
var pragma string

//go:embed sql/episode_seen_schema.sql
var episodeSeen string

func CreateTables(db *sql.DB) error {
	log := utils.GetLogger()
	log.Info("Creating database tables if not present")

	CreateTable(pragma, db) // Set pragma values, always gives error

	err := CreateTable(listItem, db)
	if err != nil {
		log.Fatal("Failed to create list_item table", "err", err)
	}
	err = CreateTable(listItemTitles, db)
	if err != nil {
		log.Fatal("Failed to create list_item_titles table", "err", err)
	}
	err = CreateTable(media, db)
	if err != nil {
		log.Fatal("Failed to create media table", "err", err)
	}
	err = CreateTable(episodeSeen, db)
	if err != nil {
		log.Fatal("Failed to create title_seen table", "err", err)
	}
	// err = CreateTable(season, db)
	// if err != nil {
	// 	log.Fatal("Failed to create seasons table", "err", err)
	// }
	err = CreateTable(sites, db)
	if err != nil {
		log.Fatal("Failed to create sites table", "err", err)
	}
	return nil
}

func CreateTable(schemaSQL string, db *sql.DB) error {
	_, err := db.Exec(schemaSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %s", err)
	}
	return nil
}
