package list

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"kultivointi-lista/db"
	"kultivointi-lista/utils"
	"strings"
)

//go:embed sql/upsert_episodes_seen.sql
var upsertEpisodesSeenSQL string

//go:embed sql/select_episodes_seen_by_item_id.sql
var selectEpisodesSeenByItemIdSQL string

//go:embed sql/select_missing_episode_seen_ids_by_item_id.sql
var selectMissingEpisodeSeenByItemIdSQL string

type EpisodeSeen struct {
	Id				string	`json:"id"`
	EpisodesSeen	int		`json:"episodesSeen"`
	SiteId			string	`json:"siteId"`
	Site 			*Site	`json:"site"`
	ItemId			string	`json:"itemId"`
}

func GetEpisodeSeenByItemID(itemId string, appCtx context.Context, tx *sql.Tx) (*EpisodeSeen, error) {
	log := utils.Logger
	episodeSeen := &EpisodeSeen{}
	itemId = strings.TrimSpace(itemId)
	if itemId == "" {
		return episodeSeen, fmt.Errorf("Missing required data item ID: %s", itemId)
	}

	mctx, err := db.MaybeCreateTx(appCtx, tx)
	if err != nil {
		return episodeSeen, err
	}
	defer mctx.MaybeCloseConn()

	row := tx.QueryRow(selectEpisodesSeenByItemIdSQL, itemId)
	err = row.Scan(&episodeSeen.Id, &episodeSeen.EpisodesSeen, &episodeSeen.SiteId, &episodeSeen.ItemId)
	if err != nil {
		log.Error("failed to scan episode seen row", "err", err)
		return episodeSeen, err
	}

	if episodeSeen.SiteId != "" {
		site, err := GetSiteById(episodeSeen.SiteId, appCtx, tx)
		if err != nil {
			return &EpisodeSeen{}, err
		}
		episodeSeen.Site = site
	}

	if err = mctx.MaybeCommit(true); err != nil {
		return &EpisodeSeen{}, err
	}

	return episodeSeen, nil
}

type EpisodesSeen []*EpisodeSeen

func (e EpisodesSeen) DeleteMissingByItemID(itemId string, appCtx context.Context, tx *sql.Tx) error {
	log := utils.GetLogger()
	currentSiteIDs := []string{}
	for _, episodeSeen := range e {
		currentSiteIDs = append(currentSiteIDs, episodeSeen.Id)
	}

	mctx, err := db.MaybeCreateTx(appCtx, tx)
	if err != nil {
		return err
	}
	defer mctx.MaybeCloseConn()

	sqlIdStr := strings.Join(currentSiteIDs, ",")
	rows, err := mctx.Tx.Query(selectMissingEpisodeSeenByItemIdSQL, itemId, sqlIdStr)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error("Failed to get missing episode seen IDs from DB", "err", err)
			return err
		}
	}

	for rows.Next() {
		episodeSeenId := ""
		err := rows.Scan(&episodeSeenId)
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				log.Error("failed to rollback changes", "err", err)
				return rbErr
			}
			return err
		}
		episodeSeen := &EpisodeSeen{
			Id: episodeSeenId,
		}
		err = episodeSeen.Delete(appCtx, tx) // should rollback here if necessary
		if err != nil {
			log.Error("failed to delete", "episodeSeen", episodeSeenId, "err", err)
			return err
		}
	}

	err = mctx.MaybeCommit(true)
	if err != nil {
		log.Error("failed to commit delete", "err", err)
		return err
	}
	return nil
}

func GetEpisodesSeenByItemID(itemId string, appCtx context.Context, tx *sql.Tx) (EpisodesSeen, error) {
	log := utils.Logger
	episodesSeen := EpisodesSeen{}

	txOnly := true
	if tx == nil {
		txOnly = false
		conn, dbCtx, err := db.GetConn(appCtx)
		if err != nil {
			return episodesSeen, err
		}
		newTx, err := conn.BeginTx(dbCtx, nil)
		if err != nil {
			return episodesSeen, err
		}
		tx = newTx
	}

	rows, err := tx.Query(selectEpisodesSeenByItemIdSQL, itemId)
	if err != nil {
		if err == sql.ErrNoRows {
			return episodesSeen, nil
		} else {
			log.Error("failed to get episodes seen rows from DB", "err", err)
		}
	}

	for rows.Next() {
		episodeSeen := &EpisodeSeen{}
		err := rows.Scan(&episodeSeen.Id, &episodeSeen.EpisodesSeen, &episodeSeen.SiteId, &episodeSeen.ItemId)
		if err != nil {
			log.Error("failed to scan episode seen row", "err", err)
			return episodesSeen, err
		}
		if episodeSeen.SiteId != "" {
			site, err := GetSiteById(episodeSeen.SiteId, appCtx, tx)
			if err != nil {
				log.Warn("failed to get site by id", "err", err)
			}
			episodeSeen.Site = site
		}
		episodesSeen = append(episodesSeen, episodeSeen)
	}

	if !txOnly {
		err := tx.Commit()
		if err != nil {
			return episodesSeen, err
		}
	}
	return episodesSeen, nil
}

func (e *EpisodeSeen) Save(appCtx context.Context, tx *sql.Tx) error {
	mctx, err := db.MaybeCreateTx(appCtx, tx)
	if err != nil {
		return err
	}
	defer mctx.MaybeCloseConn()

	if e.Id == "" {
		e.Id = db.NewID()
	}

	if e.ItemId == "" {
		tx.Rollback()
		return fmt.Errorf("Item id not set")
	}

	if e.Site != nil && e.Site.Id != "" {
		e.SiteId = e.Site.Id
		if err = e.Site.Save(appCtx, tx); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return rbErr
			}
			return err
		}
	}

	_, err = tx.Exec(upsertEpisodesSeenSQL,
		e.Id, e.EpisodesSeen, e.SiteId, e.ItemId,
		)
	if err != nil {
		return err
	}

	if err = mctx.MaybeCommit(true); err != nil {
		return err
	}

	return nil
}

func (e *EpisodeSeen) Delete(appCtx context.Context, tx *sql.Tx) error {

	mctx, err := db.MaybeCreateTx(appCtx, tx)
	if err != nil {
		return err
	}
	defer mctx.MaybeCloseConn()

	if e.Id == "" {
		e.Id = db.NewID()
	}

	_, err = mctx.Tx.Exec("DELETE FROM episode_seen WHERE id = ?", e.Id)
	if err != nil {
		return err
	}

	err = mctx.MaybeCommit(true)
	if err != nil {
		return err
	}
	return nil
}
