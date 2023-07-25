package database

import (
	_ "embed"

	"context"
	"database/sql"
	"log"
	"time"
)

var (
	//go:embed schema/tables.sql
	tablesStmt string
	//go:embed schema/triggers.sql
	triggersStmt string
)

func Init(db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, tablesStmt)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.ExecContext(ctx, triggersStmt)
	if err != nil {
		log.Fatal(err)
	}
}
