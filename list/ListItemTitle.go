package list

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"kultivointi-lista/db"
	"kultivointi-lista/utils"
	"strings"

	"github.com/charmbracelet/log"
)

//go:embed sql/upsert_list_item_title.sql
var upsertListItemTitleSQL string

//go:embed sql/select_list_item_title_by_item_id.sql
var selectListItemTitleByItemIdSQL string

//go:embed sql/select_missing_title_ids_by_item_id.sql
var selectMissingTitleIDsByItemIdSQL string

type ListItemTitle struct {
	Id				string	`json:"id"`
	Title			string	`json:"title"`
	Lang			string	`json:"lang"` // common: zh (Chinese), zh_romaji, en, jp, jp_romaji
	PrimaryTitle	bool	`json:"primaryTitle"`
	ItemId			string	`json:"itemId"`
}

func GetListItemTitleByItemID(itemId string, appCtx context.Context, tx *sql.Tx) (*ListItemTitle, error) {
	log := utils.Logger
	itemTitle := &ListItemTitle{}
	itemId = strings.TrimSpace(itemId)
	if itemId == "" {
		return itemTitle, fmt.Errorf("Missing required data item ID: %s", itemId)
	}
	txOnly := true
	if tx == nil {
		txOnly = false
		conn, dbCtx, err := db.GetConn(appCtx)
		if err != nil {
			return itemTitle, err
		}
		defer conn.Close()
		newTx, err := conn.BeginTx(dbCtx, nil)
		if err != nil {
			return itemTitle, err
		}
		tx = newTx
	}

	// SELECT id, title, lang, item_id, primary_title FROM list_item_titles WHERE item_id = ?;
	row := tx.QueryRow(selectListItemTitleByItemIdSQL, itemId)
	err := row.Scan(&itemTitle.Id, &itemTitle.Title, &itemTitle.Lang, &itemTitle.ItemId, &itemTitle.PrimaryTitle)
	if err != nil {
		log.Error("failed to scan item title row", "err", err)
		return itemTitle, err
	}

	if !txOnly {
		err := tx.Commit()
		if err != nil {
			return itemTitle, err
		}
	}
	return itemTitle, nil
}

type ItemTitles []*ListItemTitle

func (t ItemTitles) DeleteMissingByItemID(itemId string, appCtx context.Context, tx *sql.Tx) error {
	log := utils.GetLogger()
	currentTitleIDs := []string{}
	for _, title := range t {
		currentTitleIDs = append(currentTitleIDs, title.Id)
	}
	mctx, err := db.MaybeCreateTx(appCtx, tx)
	if err != nil {
		return err
	}
	defer mctx.MaybeCloseConn()
	sqlIdStr := strings.Join(currentTitleIDs, ",")
	rows, err := tx.Query(selectMissingTitleIDsByItemIdSQL, itemId, sqlIdStr)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error("Failed to get missing title IDs from DB", "err", err)
			return err
		}
	}

	for rows.Next() {
		titleId := ""
		err := rows.Scan(&titleId)
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				log.Error("failed to rollback changes", "err", err)
				return rbErr
			}
			return err
		}
		title := &ListItemTitle{
			Id: titleId,
		}
		err = title.Delete(appCtx, tx) // should rollback here if necessary
		if err != nil {
			log.Error("failed to delete", "titleId", titleId, "err", err)
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

type ItemTitlesSearchOptions struct {
	OrderDirection string `json:"orderDirection"` // ASC / DESC
}
func (l *ItemTitles) Search(appCtx context.Context, query string, options *ItemTitlesSearchOptions) error {
	log := utils.GetLogger()
	query = strings.TrimSpace(query)
	query = strings.ReplaceAll(query, " ", "%")
	query = fmt.Sprintf("%%%s%%", query)
	searchSQL := fmt.Sprintf("SELECT id, title, lang, item_id, primary_title FROM list_item_titles WHERE title LIKE ? ORDER BY title %s", options.OrderDirection)

	orderDefault := "DESC"
	if(options == nil) {
		options = &ItemTitlesSearchOptions{
			OrderDirection: orderDefault,
		}
	}
	if(options.OrderDirection == "" || !(options.OrderDirection == "ASC" || options.OrderDirection == "DESC" )) {
		log.Warn("Invalid order direction, using default", "OrderDirection", options.OrderDirection)
		options.OrderDirection = orderDefault
	}
	conn, dbCtx, err := db.GetConn(appCtx)
	if err != nil {
		log.Error("Failed to get connection", "err", err);
	}
	defer conn.Close()
	rows, err := conn.QueryContext(dbCtx, searchSQL, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			return err
		}
	}
	for rows.Next() {
		item := &ListItemTitle{}
		err := rows.Scan(&item.Id, &item.Title, &item.Lang, &item.ItemId, &item.PrimaryTitle)
		if err != nil {
			log.Error("failed to scan row")
			continue
		}
		*l = append(*l, item)
	}
	return nil
}

func GetListItemTitlesByItemID(itemId string, appCtx context.Context, tx *sql.Tx) (ItemTitles, error) {
	log := utils.Logger
	titles := ItemTitles{}

	txOnly := true
	if tx == nil {
		txOnly = false
		conn, dbCtx, err := db.GetConn(appCtx)
		if err != nil {
			return titles, err
		}
		defer conn.Close()
		newTx, err := conn.BeginTx(dbCtx, nil)
		if err != nil {
			return titles, err
		}
		tx = newTx
	}

	rows, err := tx.Query(selectListItemTitleByItemIdSQL, itemId)
	if err != nil {
		if err == sql.ErrNoRows {
			return titles, nil
		} else {
			log.Error("failed to get item title rows from DB", "err", err)
		}
	}

	for rows.Next() {
		itemTitle := &ListItemTitle{}
		// SELECT id, title, lang, item_id, primary_title FROM list_item_titles WHERE item_id = ?;
		err := rows.Scan(&itemTitle.Id, &itemTitle.Title, &itemTitle.Lang, &itemTitle.ItemId, &itemTitle.PrimaryTitle)
		if err != nil {
			log.Error("failed to scan item title row", "err", err)
			return titles, err
		}
		titles = append(titles, itemTitle)
	}

	if !txOnly {
		err := tx.Commit()
		if err != nil {
			return titles, err
		}
	}
	return titles, nil
}

func (t *ItemTitles) GetPrimary() (*ListItemTitle, error) {
	if len(*t) < 1 {
		return &ListItemTitle{}, fmt.Errorf("No titles found")
	}
	for _, title := range *t {
		if title.PrimaryTitle {
			return title, nil
		}
	}
	return &ListItemTitle{}, fmt.Errorf("Primary title not found")
}

func (t *ListItemTitle) Save(appCtx context.Context, tx *sql.Tx) error {

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

	if t.Id == "" {
		t.Id = db.NewID()
	}
	if t.Lang == "" {
		t.Lang = "zh_romanji"
	}
	if t.ItemId == "" {
		if !txOnly {
			err := tx.Rollback()
			if err != nil {
				return err
			}
		}
		return fmt.Errorf("No item ID set")
	}

	_, err := tx.Exec(upsertListItemTitleSQL,
		t.Id, t.Title, t.Lang, t.ItemId, t.PrimaryTitle,
		)
	if err != nil {
		return err
	}

	if !txOnly {
		err := tx.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *ListItemTitle) Delete(appCtx context.Context, tx *sql.Tx) error {
	mctx, err := db.MaybeCreateTx(appCtx, tx)
	if err != nil {
		return err
	}
	// txOnly := true
	// if tx == nil {
	// 	txOnly = false
	// 	conn, dbCtx, err := db.GetConn(appCtx)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	defer conn.Close()
	// 	newTx, err := conn.BeginTx(dbCtx, nil)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	tx = newTx
	// }

	if t.Id == "" {
		if !mctx.TxOnly {
			if rbErr := mctx.Tx.Rollback(); rbErr != nil {
				db.MaybeLogError(rbErr, "No ID set failed to rollback", "err", rbErr)
				return fmt.Errorf("No ID set")
			}
		}
		return fmt.Errorf("No ID set")
	}

	_, err = mctx.Tx.Exec("DELETE FROM list_item_titles WHERE id = ?", t.Id)
	if err != nil {
		if !mctx.TxOnly {
			if rbErr := mctx.Tx.Rollback(); rbErr != nil {
				db.MaybeLogError(rbErr, "failed to rollback", "err", rbErr)
			}
		}
		log.Error("Failed to delete list item title", "err", err)
		return err
	}

	// if !txOnly {
	// 	err := tx.Commit()
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	err = mctx.MaybeCommit(true)
	if err != nil {
		if !mctx.TxOnly {
			db.MaybeLogError(err, "Failed to delete list item", "title", t.Title)
		} else {
			db.MaybeLogError(err, "Failed to delete list item in transaction", "title", t.Title)
		}
	}
	return nil
}
