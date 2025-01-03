package db

import (
	"context"
	"database/sql"
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"

	// _ "github.com/tursodatabase/go-libsql"
	"kultivointi-lista/utils"

	_ "github.com/mattn/go-sqlite3"
)

// Generate a ID for DB, configuration here
func NewID() string {
	log := utils.GetLogger()
	id, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 15)
	if err != nil {
		log.Fatal("failed to generate nano id")
	}
	return id
}

func GetConn(appCtx context.Context) (conn *sql.Conn, connCtx context.Context, err error) {
	// log := utils.GetLogger()
	dbPool := appCtx.Value("dbPool").(*sql.DB)
	if dbPool == nil {
		return nil, nil, fmt.Errorf("Failed to get DB Pool from context")
	}
	fmt.Println(dbPool)
	dbCtx := context.Background()
	conn, err = dbPool.Conn(dbCtx)

	return conn, dbCtx, err
}

func TursoConnect(dbFile string) (*sql.DB, error) {
	log := utils.GetLogger()
	connStr := fmt.Sprintf("file:%s", dbFile)
	log.Info("DB Connection string", "value", connStr)
	// sqlite3
	db, err := sql.Open("sqlite3", connStr)
	// db, err := sql.Open("libsql", connStr)
	if err != nil {
		errorMsg := fmt.Sprintf("TursoConnect: failed to open db")
		log.Error(errorMsg, "err", err)
		return nil, fmt.Errorf("%s: err: %s", errorMsg, err)
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	db.SetConnMaxIdleTime(2000)
	// _, err = db.Exec("PRAGMA busy_timeout = 30000;")
	// 			if err != nil {
	// 				log.Error("Failed to proagma modes", "err", err);
	// 				return nil, err
	// 			}
	return db, nil
}

type MaybeCreateTxResponse struct {
	TxOnly bool
	Tx *sql.Tx
	Conn *sql.Conn
	DbCtx context.Context
}

/*
// USE:
mctx, _ := MaybeCreateTx()
defer mctx.MaybeCloseConn()
// logic here...
if err = mctx.MaybeCommit(true); err != nil {
	return err
}
*/
func MaybeCreateTx(appCtx context.Context, tx *sql.Tx) (*MaybeCreateTxResponse, error) {
	res := &MaybeCreateTxResponse{
		TxOnly: true,
	}
	if tx == nil {
		res.TxOnly = false
		conn, dbCtx, err := GetConn(appCtx)
		if err != nil {
			return &MaybeCreateTxResponse{}, err
		}
		res.Conn = conn
		res.DbCtx = dbCtx
		newTx, err := conn.BeginTx(dbCtx, nil)
		if err != nil {
			MaybeLogError(err, "Failed to create transaction")
			if rbErr := tx.Rollback(); rbErr != nil {
				MaybeLogError(rbErr, "failed to rollback error", "err", rbErr)
			}
			return &MaybeCreateTxResponse{}, err
		}
		tx = newTx
	}
	res.Tx = tx
	return res, nil
}

// if err = mctx.MaybeCommit(true); err != nil {
// 		return err
// 	}
func (t *MaybeCreateTxResponse) MaybeCommit(rollbackOnError bool) error {
	if !t.TxOnly {
		err := t.Tx.Commit()
		if err != nil {
			if rollbackOnError {
				rbErr := t.Tx.Rollback()
				if rbErr != nil {
					MaybeLogError(rbErr, "failed to rollback changes", "err", err)
					return rbErr
				}
			}
			MaybeLogError(err, "failed to commit changes to DB", "err", err)
			return err
		}
	}
	return nil
}

// USE after mctx, _ := MaybeCreateTx()
// defer mctx.MaybeCloseConn()
func (t *MaybeCreateTxResponse) MaybeCloseConn() {
	if t.TxOnly == false && t.Conn != nil {
		t.Conn.Close()
	}
}

func MaybeLogError(originalError error, msg interface{}, keyvals ...interface{}) {
	log := utils.GetLogger()
	if originalError != sql.ErrTxDone {
		log.Error(msg, keyvals...)
	}
}
