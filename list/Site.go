package list
import (
	"database/sql"
	"kultivointi-lista/db"
	"kultivointi-lista/utils"
	"context"
	"strings"
	_ "embed"
)

//go:embed sql/upsert_site.sql
var upsertSiteSQL string

//go:embed sql/select_missing_site_ids_by_item_id.sql
var selectMissingTitleIDsSQL string

//go:embed sql/select_site_by_id.sql
var selectSiteByIdSQL string

//go:embed sql/select_all_sites.sql
var selectAllSitesSQL string

type Site struct {
	Id					string	`json:"id"`
	Url					string	`json:"url"`
	DomainBase			string	`json:"domainBase"`
	DomainTopLevel		string	`json:"domainTopLevel"` // .com .net .io etc
	DomainProtocol		string	`json:"domainProtocol"` // http / https
	EpisodeTemplate		string	`json:"episodeTemplate"` // Simple string replaceable DSL for generating URLS for episodes
	MainPageTemplate	string	`json:"mainPageTempate"` // Simple DSL for generating for movie/serie URLs
}

func GetSiteById(siteId string, appCtx context.Context, tx *sql.Tx) (*Site, error) {
	log := utils.GetLogger()
	mctx, err := db.MaybeCreateTx(appCtx, tx)
	if err != nil {
		log.Error("failed to create TX", "err", err)
	}
	defer mctx.MaybeCloseConn()
	row := mctx.Tx.QueryRow(selectSiteByIdSQL, siteId)
	site := &Site{}
	err = row.Scan(
		&site.Id,
		&site.Url,
		&site.DomainBase,
		&site.DomainTopLevel,
		&site.DomainProtocol,
		&site.EpisodeTemplate,
		&site.MainPageTemplate,
	)
	if err != nil {
		log.Error("failed to get site from DB", "err", err)
		rbErr := mctx.Tx.Rollback()
		if rbErr != nil {
			return &Site{}, rbErr
		}
		return &Site{}, err
	}

	err = mctx.MaybeCommit(true)
	if err != nil {
		log.Error("failed to commit changes to DB", "err", err)
	}
	return site, nil
}
type Sites []*Site

func GetAllSites(appCtx context.Context, tx *sql.Tx) (Sites, error) {
	log := utils.GetLogger()
	mctx, err := db.MaybeCreateTx(appCtx, tx)
	if err != nil {
		log.Error("failed to create TX", "err", err)
		return Sites{}, err
	}
	defer mctx.MaybeCloseConn()

	rows, err := mctx.Tx.Query(selectAllSitesSQL)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error("failed to get sites from DB", "err", err)
			return Sites{}, err
		} else {
			return Sites{}, nil
		}
	}

	allSites := Sites{}
	for rows.Next() {
		site := &Site{}
		err = rows.Scan(
			&site.Id,
			&site.Url,
			&site.DomainBase,
			&site.DomainTopLevel,
			&site.DomainProtocol,
			&site.EpisodeTemplate,
			&site.MainPageTemplate,
		)
		if err != nil {
			log.Error("failed to scan site row", "err", err)
			continue
		}
		allSites = append(allSites, site)
	}

	if err = mctx.MaybeCommit(true); err != nil {
		return Sites{}, err
	}
	return allSites, nil
}
func (s Sites) DeleteMissingByItemID(itemId string, appCtx context.Context, tx *sql.Tx) error {
	log := utils.GetLogger()
	currentSiteIDs := []string{}
	for _, site := range s {
		currentSiteIDs = append(currentSiteIDs, site.Id)
	}

	mctx, err := db.MaybeCreateTx(appCtx, tx)
	if err != nil {
		return err
	}
	defer mctx.MaybeCloseConn()

	sqlIdStr := strings.Join(currentSiteIDs, ",")
	rows, err := tx.Query(selectMissingTitleIDsByItemIdSQL, itemId, sqlIdStr)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error("Failed to get missing site IDs from DB", "err", err)
			return err
		}
	}

	for rows.Next() {
		siteId := ""
		err := rows.Scan(&siteId)
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				log.Error("failed to rollback changes", "err", err)
				return rbErr
			}
			return err
		}
		site := &Site{
			Id: siteId,
		}
		err = site.Delete(appCtx, tx) // should rollback here if necessary
		if err != nil {
			log.Error("failed to delete", "siteId", siteId, "err", err)
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

func (t *Site) Save(appCtx context.Context, tx *sql.Tx) error {

	mctx, err := db.MaybeCreateTx(appCtx, tx)
	if err != nil {
		return err
	}
	defer mctx.MaybeCloseConn()

	if t.Id == "" {
		t.Id = db.NewID()
	}

	_, err = mctx.Tx.Exec(upsertSiteSQL,
		t.Id, t.Url, t.DomainBase, t.DomainTopLevel, t.DomainProtocol, t.EpisodeTemplate, t.MainPageTemplate,
		)
	if err != nil {
		return err
	}

	err = mctx.MaybeCommit(true)
	if err != nil {
		return err
	}
	return nil
}

func (t *Site) Delete(appCtx context.Context, tx *sql.Tx) error {

	mctx, err := db.MaybeCreateTx(appCtx, tx)
	if err != nil {
		return err
	}
	defer mctx.MaybeCloseConn()

	if t.Id == "" {
		t.Id = db.NewID()
	}

	_, err = tx.Exec("DELETE FROM sites WHERE id = ?", t.Id)
	if err != nil {
		return err
	}

	err = mctx.MaybeCommit(true)
	if err != nil {
		return err
	}
	return nil
}
