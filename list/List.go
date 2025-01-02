package list

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"kultivointi-lista/db"
	"kultivointi-lista/utils"
	"time"

	"github.com/charmbracelet/log"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed sql/select_base_item_ids.sql
var selectBaseItemIDsSQL string

type List struct {
	// MainList []*ListItem `json:"mainList"`
	// Filtered []*ListItem `json:"filtered"`
	getMainListCh chan mainListChResult
	setListItemTitleCh chan setListItemTitleChResult
	ctx context.Context
}
type setListItemTitleChResult struct {
	listItem *ListItem
	err error
}
type mainListChResult struct {
	list []*ListItem
	err error
}
func NewList(appCtx context.Context) *List {
	return &List{
		getMainListCh: make(chan mainListChResult),
		setListItemTitleCh: make(chan setListItemTitleChResult),
		ctx: appCtx,
	}
}

// Go routine that handles lists
func (l *List) ListChannels() {
	for {
		select {
			case _, ok := <-l.getMainListCh:
			if !ok {
				l.getMainListCh <- mainListChResult{
					list: []*ListItem{},
					err: fmt.Errorf("channel ok false"),
				}
			}
			listItemIDs, err := l.GetBaseItemIDs()
			list, err := GetListByIDs(l.ctx, listItemIDs)
			l.getMainListCh <- mainListChResult{
				list: list,
				err: err,
			}
			default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

type ListSearchOptions struct {
	SearchFields []string `json:"searchFields"` // title
	OrderField string `json:"orderField"` // title, modified
	OrderDirection string `json:"orderDirection"` // ASC / DESC
}
func (l *List) Search(query string, options *ListSearchOptions) ([]*ListItem, error) {
	defaultSearchFields := []string{ "title" }
	defaultOrderField := "title"
	defaultOrderDirection := "DESC"
	if(options == nil) {
		options = &ListSearchOptions{
			SearchFields: defaultSearchFields,
			OrderField: defaultOrderField,
			OrderDirection: defaultOrderDirection,
		}
	}
	if len(options.SearchFields) == 0 {
		options.SearchFields = defaultSearchFields
	}
	if options.OrderField == "" {
		options.OrderField = defaultOrderField
	}
	if options.OrderDirection == "" || !(options.OrderDirection == "ASC" || options.OrderDirection == "DESC") {
		options.OrderDirection = defaultOrderDirection
	}
	titles := &ItemTitles{}
	err := titles.Search(l.ctx, query, &ItemTitlesSearchOptions{
		OrderDirection: options.OrderDirection,
	})
	if err != nil {
		log.Error("failed to search titles", "err", err);
		return []*ListItem{}, fmt.Errorf("failed to search titles")
	}
	titleItemIDs := []string{}
	for _, title := range *titles {
		titleItemIDs = append(titleItemIDs, title.ItemId)
	}
	list := NewList(l.ctx)
	listItems := list.GetListItemsByIDs(titleItemIDs)
	return listItems, nil
}
// Push change event to frontend if it isn't same as compareTitle (if compareTitle == nil always push)
func ListItemTitlePushEvent(appCtx context.Context, newTitle *ListItemTitle, compareTitle *ListItemTitle) {
	hasChanges := false
	if compareTitle == nil {
		hasChanges = true
	}
	if newTitle.Id != compareTitle.Id {
		hasChanges = true
	}
	if newTitle.ItemId != compareTitle.ItemId {
		hasChanges = true
	}
	if newTitle.Lang != compareTitle.Lang {
		hasChanges = true
	}
	if newTitle.Title != compareTitle.Title {
		hasChanges = true
	}
	if newTitle.PrimaryTitle != compareTitle.PrimaryTitle {
		hasChanges = true
	}
	if hasChanges {
		runtime.EventsEmit(appCtx, "list_item_title_push", newTitle)
	}
}

func (l *List) NewDbID() string {
	return db.NewID()
}

func (l *List) GetMainList() ([]*ListItem, error) {
	mainListRes := mainListChResult{}
	l.getMainListCh<-mainListRes
	mainListRes, ok := <-l.getMainListCh
	if !ok || mainListRes.err != nil {
		return mainListRes.list, fmt.Errorf("failed to get main res from channel, ok: %v, err: %v",  ok, mainListRes.err)
	}
	return mainListRes.list, nil
}

func (l *List) SaveListItem(item *ListItem) bool {
	log := utils.GetLogger()
	utils.PrettyPrint(item);
	conn, dbCtx, err := db.GetConn(l.ctx)
	if err != nil {
		log.Error("failed to aquire new connection", "err", err)
		return false
	}
	defer conn.Close()
	tx, err := conn.BeginTx(dbCtx, nil)
	if err != nil {
		log.Error("failed to begin transaction", "err", err)
		return false
	}
	err = item.Save(l.ctx, tx)
	if err != nil {
		log.Error("Failed to save list item", "err", err);
		return false
	}
	err = tx.Commit()
	if err != nil {
		log.Error("failed to commit list item save", "err", err);
		return false;
	}
	return true
}

func (l *List) DeleteListItem(itemId string) bool {
	listItem := &ListItem{ Id: itemId }
	err := listItem.Delete(l.ctx, nil)
	if err != nil {
		return false
	}
	return true
}

func (l *List) GetListItemsByIDs(itemIDs []string) []*ListItem {
	log := utils.GetLogger()
	items, err := GetListByIDs(l.ctx, itemIDs)
	if err != nil {
		log.Error("Failed to get items from DB", "err", err)
		return []*ListItem{}
	}
	return items
}

func (l *List) SetContext(appCtx context.Context) {
	l.ctx = appCtx
}

func GetListByIDs(appCtx context.Context, IDs []string) ([]*ListItem, error) {
	log := utils.Logger
	conn, dbCtx, err := db.GetConn(appCtx)
	list := []*ListItem{}
	if err != nil {
		return list, err
	}
	tx, err := conn.BeginTx(dbCtx, nil)
	if err != nil {
		return list, err
	}
	for _, itemId := range IDs {
		// fmt.Println("ID", itemId)
		item, err := GetListItemByID(itemId, appCtx, tx)
		if err != nil {
			log.Error("Failed to get list item", "err", err)
		}
		list = append(list, item)
	}
	return list, nil
}

func (l *List) GetBaseItemIDs() ([]string, error) {
	log := utils.Logger
	log.Info("Get base item IDs")
	conn, dbCtx, err := db.GetConn(l.ctx)
	if err != nil {
		return []string{}, fmt.Errorf("failed to get connection from pool")
	}
	rows, err := conn.QueryContext(dbCtx, selectBaseItemIDsSQL)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info("No rows")
			return []string{}, nil
		} else {
			log.Error("failed to get Base item IDs", "err", err)
			return []string{}, fmt.Errorf("failed to get Base item IDs")
		}
	}

	ids := []string{}
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil || id == "" {
			log.Error("Failed to get ID from DB row", "err", err)
			return []string{}, fmt.Errorf("failed to get id from DB row")
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (l *List) GetSites() (Sites, error) {
	log := utils.GetLogger()
	sites, err := GetAllSites(l.ctx, nil)
	if err != nil {
		log.Error("Failed to get sites", "err", err)
		return Sites{}, err
	}
	return sites, nil
}

func (l *List) SaveSite(site *Site) error {
	utils.PrettyPrint(site)
	log := utils.GetLogger()
	if site == nil {
		log.Error("no site to save", "site", site)
		return fmt.Errorf("no site to save")
	}
	err := site.Save(l.ctx, nil)
	if err != nil {
		log.Error("failed to save site", "err", err)
		return err
	}

	return nil
}
