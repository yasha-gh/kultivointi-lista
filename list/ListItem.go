package list

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"kultivointi-lista/db"
	"kultivointi-lista/utils"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

//go:embed sql/upsert_list_item.sql
var upsertListItemSQL string

//go:embed sql/select_list_item_by_id.sql
var selectListItemByIdSQL string

type ListItem struct {
	Id					string			`json:"id"`
	TitleId				string			`json:"titleId"`
	Title				*ListItemTitle	`json:"title"` // pointer to titles -> PrimaryTitle == true
	Titles				ItemTitles		`json:"titles"`
	Type				string			`json:"type"` // base, season
	BroadcastType		string			`json:"broadcastType"` // TV, OVA, ONA, Movie
	EpisodesTotal		int				`json:"episodesTotal"`
	EpisodesSeen		int				`json:"episodesSeen"`
	Ongoing				bool			`json:"ongoing"`
	SeasonNum			int				`json:"seasonNum"`
	Seasons				Seasons			`json:"seasons"`
	ParentItemId		string			`json:"parentItemId"` // list item ID if type is season
	EpisodesSeenOn		EpisodesSeen	`json:"episodesSeenOn"` // Sites where episodes seen
	ThumbnailImageId	string			`json:"thubmnailImageId"`
	ModifiedAt			time.Time		`json:"modifiedAt"`
	CreatedAt			time.Time		`json:"createdAt"`
	//TODO thumbnail_image_id / media
}
type Seasons []*ListItem

func GetListItemByID(id string, appCtx context.Context, tx *sql.Tx) (*ListItem, error) {
	log := utils.Logger
	item := &ListItem{}
	id = strings.TrimSpace(id)
	if id == "" {
		return item, fmt.Errorf("Missing required data ID: %s", id)
	}
	txOnly := true
	if tx == nil {
		txOnly = false
		conn, dbCtx, err := db.GetConn(appCtx)
		if err != nil {
			return item, err
		}
		defer conn.Close()
		newTx, err := conn.BeginTx(dbCtx, nil)
		if err != nil {
			return item, err
		}
		tx = newTx
	}

	// SELECT id, title_id, type, broadcast_type, thumbnail_image_id, ongoing, episodes_total, episodes_seen, parent_item_id, season_num FROM list_items WHERE id = ?;
	row := tx.QueryRow(selectListItemByIdSQL, id)
	err := row.Scan(
		&item.Id,
		&item.TitleId,
		&item.Type,
		&item.BroadcastType,
		&item.ThumbnailImageId,
		&item.Ongoing,
		&item.EpisodesTotal,
		&item.EpisodesSeen,
		&item.ParentItemId,
		&item.SeasonNum,
		&item.ModifiedAt,
		&item.CreatedAt,
	)
	if err != nil {
		// fmt.Println("item id", id);
		log.Error("failed to scan item row", "err", err)
		return item, err
	}

	titles, err := GetListItemTitlesByItemID(item.Id, appCtx, tx)
	if err != nil {
		log.Error("Failed to get titles")
		rbErr := tx.Rollback()
		if rbErr != nil {
			log.Error("failed to rollback TX", "err", err)
		}
		return item, err
	}

	if len(titles) < 1 {
		log.Error("No titles found")
		rbErr := tx.Rollback()
		if rbErr != nil {
			log.Error("failed to rollback TX", "err", err)
		}
		return item, err
	}

	primaryTitle, err := titles.GetPrimary()
	if err != nil {
		log.Error("failed to get primary title", "err", err)
		return item, err
	}

	item.Title = primaryTitle
	item.Titles = titles

	episodesSeen, err := GetEpisodesSeenByItemID(item.Id, appCtx, tx)
	if err != nil {
		log.Error("Failed to get episodes seen")
		rbErr := tx.Rollback()
		if rbErr != nil {
			log.Error("failed to rollback TX", "err", err)
		}
		return item, err
	}

	item.EpisodesSeenOn = episodesSeen

	if !txOnly {
		err := tx.Commit()
		if err != nil {
			return item, err
		}
	}
	return item, nil
}

func (l *ListItem) Save(appCtx context.Context, tx *sql.Tx) error {
	txOnly := true
	if tx == nil {
		txOnly = false
		conn, dbCtx, err := db.GetConn(appCtx)
		if err != nil {
			return err
		}
		defer conn.Close()
		newTx, err := conn.BeginTx(dbCtx, nil)
		if err != nil {
			return err
		}
		tx = newTx
	}
	if l.Id == "" {
		l.Id = db.NewID()
	}

	if l.Type == "" {
		l.Type = "base"
	}

	if l.BroadcastType == "" {
		l.BroadcastType = "ONA"
	}

	err := l.Titles.DeleteMissingByItemID(l.Id, appCtx, tx)
	if err != nil {
		log.Error("Failed to delete missing titles", "err", err)
		return err
	}

	for _, title := range l.Titles {
		if title.ItemId == ""  {
			title.ItemId = l.Id
		}
		err := title.Save(appCtx, tx)
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				return err
			}
			return err
		}
		// l.Titles[i] = *tp
	}

	if l.Title == nil || l.TitleId == "" {
		// titles := &l.Titles
		primary, err := l.Titles.GetPrimary()
		if err != nil {
			txErr := tx.Rollback()
			if txErr != nil {
				return txErr
			}
			return err
		}
		l.Title = primary
		l.TitleId = l.Title.Id
	}

	err = l.EpisodesSeenOn.DeleteMissingByItemID(l.Id, appCtx, tx)
	if err != nil {
		log.Error("failed to delete missing episode seen on titles", "err", err)
	}

	for _, seenOn := range l.EpisodesSeenOn {
		if seenOn.ItemId == "" {
			seenOn.ItemId = l.Id
		}
		err := seenOn.Save(appCtx, tx)
		if err != nil {
			txErr := tx.Rollback()
			if txErr != nil {
				return err
			}
			return err
		}
		// l.EpisodesSeenOn[i] = *seenOnP
	}
	modifiedAt := time.Now()
	_, err = tx.Exec(upsertListItemSQL,
		l.Id,
		l.TitleId,
		l.Type,
		l.BroadcastType,
		l.ThumbnailImageId,
		l.Ongoing,
		l.EpisodesTotal,
		l.EpisodesSeen,
		l.ParentItemId,
		l.SeasonNum,
		modifiedAt,
		)
	if err != nil {
		return err
	}

	if !txOnly {
		log.Info("List item saved", "title", l.Title.Title)
		err := tx.Commit()
		if err != nil {
			return err
		}
	} else {
		log.Info("List item save success in transaction", "title", l.Title.Title)
	}
	return nil
}

func (l *ListItem) Delete(appCtx context.Context, tx *sql.Tx) error {
	log := utils.GetLogger()
	log.Info("Delete list item", "title", l.Title)
	if l.Id == "" {
		return fmt.Errorf("No ID on list item")
	}
	mctx, err := db.MaybeCreateTx(appCtx, tx)
	if err != nil {
		log.Error("Failed to get tx", "err", err)
		return err
	}
	defer mctx.MaybeCloseConn()

	for _, title := range l.Titles {
		if err := title.Delete(appCtx, mctx.Tx); err != nil {
			if rbErr := mctx.Tx.Rollback(); rbErr != nil {
				log.Error("Failed to rollback", "err", rbErr)
			}
			return err
		}
	}

	for _, epSeen := range l.EpisodesSeenOn {
		if err := epSeen.Delete(appCtx, mctx.Tx); err != nil {
			if rbErr := mctx.Tx.Rollback(); rbErr != nil {
				db.MaybeLogError(rbErr, "Failed to rollback", "err", rbErr)
				return rbErr
			}
			return err
		}
	}

	_, err = mctx.Tx.Exec("DELETE FROM list_items WHERE id = ?", l.Id)
	if err != nil {
		log.Error("Failed to delete list item", "itemID", l.Id, "err", err)
		if rbErr := mctx.Tx.Rollback(); rbErr != nil {
			db.MaybeLogError(rbErr, "Failed to rollback", "err", rbErr)
			return rbErr
		}
		return err
	}
	deleteTitle := ""
	if l.Title != nil {
		deleteTitle = l.Title.Title
	}
	if err = mctx.MaybeCommit(true); err != nil {
		log.Info("deleted list item ", "title", deleteTitle)
		return err
	} else {
		log.Info("delete list item success in transaction", "title", deleteTitle)
	}
	return nil
}
