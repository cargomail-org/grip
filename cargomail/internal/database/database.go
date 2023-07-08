package database

import (
	"context"
	"database/sql"
	"log"
	"time"
)

func Init(db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS user (
		id				INTEGER PRIMARY KEY,
		username		TEXT UNIQUE NOT NULL,
		password_hash	TEXT NOT NULL,
		created_at		TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS session (
		hash 			BLOB PRIMARY KEY,
		user_id 		INTEGER NOT NULL REFERENCES user ON DELETE CASCADE,
		expiry 			TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		scope 			TEXT NOT NULL
 	);
	CREATE TABLE IF NOT EXISTS resources (
		id				INTEGER PRIMARY KEY,
		user_id 		INTEGER NOT NULL REFERENCES user ON DELETE CASCADE,
		uuid			TEXT UNIQUE,
		name			TEXT NOT NULL,
		path			TEXT NOT NULL,
		size			INTEGER NOT NULL,
		content_type	TEXT NOT NULL,
		created_at		TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		log.Fatal(err)
	}
}
