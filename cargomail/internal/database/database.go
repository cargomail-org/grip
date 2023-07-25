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
	------------------------------tables-----------------------------

	CREATE TABLE IF NOT EXISTS contact_timeline_seq (
		last_timeline_id integer(8) NOT NULL
	);
	CREATE TABLE IF NOT EXISTS contact_history_seq (
		last_history_id integer(8) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS user (
		id				INTEGER PRIMARY KEY,
		username		TEXT NOT NULL UNIQUE,
		password_hash	TEXT NOT NULL,
		firstname		TEXT DEFAULT "",
		lastname		TEXT DEFAULT "",
		created_at		TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS session (
		hash 			BLOB PRIMARY KEY,
		user_id 		INTEGER NOT NULL REFERENCES user ON DELETE CASCADE,
		expiry 			TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		scope 			TEXT NOT NULL
 	);
	CREATE TABLE IF NOT EXISTS file (
		id				INTEGER PRIMARY KEY,
		user_id 		INTEGER NOT NULL REFERENCES user ON DELETE CASCADE,
		uuid			TEXT NOT NULL UNIQUE,
		hash 			BLOB NOT NULL,
		name			TEXT NOT NULL,
		path			TEXT NOT NULL,
		size			INTEGER NOT NULL,
		content_type	TEXT NOT NULL,
		created_at		TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS contact (
		id				INTEGER PRIMARY KEY,
		user_id 		INTEGER NOT NULL REFERENCES user ON DELETE CASCADE,
		uuid			TEXT NOT NULL UNIQUE,
		email_address   TEXT,
		firstname		TEXT,
		lastname		TEXT,
		created_at		TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		modified_at		TIMESTAMP,
		timeline_id		INTEGER(8) NOT NULL DEFAULT 0,
		history_id 		INTEGER(8) NOT NULL DEFAULT 0,
		last_stmt  		INTEGER(2) NOT NULL DEFAULT 0 -- 0-inserted, 1-updated, 2-trashed
	);

	------------------------------indexes----------------------------

	CREATE INDEX IF NOT EXISTS idx_file_hash ON file(hash);

	CREATE INDEX IF NOT EXISTS idx_contact_timeline_id ON contact (timeline_id);
	CREATE INDEX IF NOT EXISTS idx_contact_history_id ON contact (history_id);
	CREATE INDEX IF NOT EXISTS idx_contact_last_stmt ON contact (last_stmt);

	------------------------------initialization---------------------

	INSERT INTO contact_timeline_seq(last_timeline_id) 
		SELECT 1
		WHERE NOT EXISTS (SELECT 1 from contact_timeline_seq);

	INSERT INTO contact_history_seq(last_history_id) 
		SELECT 1
		WHERE NOT EXISTS (SELECT 1 from contact_history_seq);

	------------------------------contacts triggers---------------------------		

	CREATE TRIGGER IF NOT EXISTS contact_after_insert
		AFTER INSERT
		ON contact
		FOR EACH ROW
	BEGIN
		UPDATE contact_timeline_seq SET last_timeline_id = (last_timeline_id + 1);
		UPDATE contact_history_seq SET last_history_id = (last_history_id + 1);
		UPDATE contact
		SET timeline_id = (SELECT last_timeline_id FROM contact_timeline_seq),
			history_id  = (SELECT last_history_id FROM contact_history_seq),
			last_stmt   = 0
		WHERE id = new.id;
	END;
	
	CREATE TRIGGER IF NOT EXISTS contact_before_update
		BEFORE UPDATE OF
			id,
			uuid
		ON contact
		FOR EACH ROW
	BEGIN
		SELECT RAISE(ABORT, 'Update not allowed');
	END;
	
	CREATE TRIGGER IF NOT EXISTS contact_after_update
		AFTER UPDATE OF
			email_address,
			firstname,
			lastname
		ON contact
		FOR EACH ROW
	BEGIN
		UPDATE contact_timeline_seq SET last_timeline_id = (last_timeline_id + 1);
		UPDATE contact_history_seq SET last_history_id = (last_history_id + 1);
		UPDATE contact
		SET timeline_id = (SELECT last_timeline_id FROM contact_timeline_seq),
			history_id  = (SELECT last_history_id FROM contact_history_seq),
			last_stmt   = 1,
			modified_at = CURRENT_TIMESTAMP
		WHERE id = old.id;
	END;
	
	-- Trashed
	CREATE TRIGGER IF NOT EXISTS contact_before_trashed
		BEFORE UPDATE OF
			last_stmt
		ON contact
		FOR EACH ROW
	BEGIN
		SELECT RAISE(ABORT, 'Update "last_stmt" not allowed')
		WHERE (new.last_stmt < 0 OR new.last_stmt > 2)
		   OR (old.last_stmt = 2 AND new.last_stmt != 0); -- Untrash = trashed (2) -> inserted (0)
	END;
	
	CREATE TRIGGER IF NOT EXISTS contact_after_trashed
		AFTER UPDATE OF
			last_stmt
		ON contact
		FOR EACH ROW
		WHEN new.last_stmt = 2
	BEGIN
		UPDATE contact_history_seq SET last_history_id = (last_history_id + 1);
		UPDATE contact
		SET history_id  = (SELECT last_history_id FROM contact_history_seq),
			last_stmt   = new.last_stmt,
			modified_at = CURRENT_TIMESTAMP
		WHERE id = old.id;
	END;
	
	`)
	if err != nil {
		log.Fatal(err)
	}
}
