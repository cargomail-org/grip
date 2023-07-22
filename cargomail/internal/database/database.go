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

	CREATE TABLE IF NOT EXISTS timeline_seq (
			num integer(8) NOT NULL
	);
	CREATE TABLE IF NOT EXISTS history_seq (
		num integer(8) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS user (
		id				INTEGER PRIMARY KEY,
		username		TEXT NOT NULL UNIQUE,
		password_hash	TEXT NOT NULL,
		firstname		TEXT,
		lastname		TEXT,
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
		timeline_id		INTEGER(8) NOT NULL DEFAULT 0,
		history_id 		INTEGER(8) NOT NULL DEFAULT 0,
		last_stmt  		INTEGER(2) NOT NULL DEFAULT 0, -- 0-insert, 1-update, 2-mark for delete
		timestamp		TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	------------------------------indexes----------------------------

	CREATE INDEX IF NOT EXISTS idx_file_hash ON file(hash);

	CREATE INDEX IF NOT EXISTS idx_contact_timeline_id ON contact (timeline_id);
	CREATE INDEX IF NOT EXISTS idx_contact_history_id ON contact (history_id);
	CREATE INDEX IF NOT EXISTS idx_contact_last_stmt ON contact (last_stmt);

	------------------------------initialization---------------------

	INSERT INTO timeline_seq(num) 
		SELECT 1
		WHERE NOT EXISTS (SELECT 1 from timeline_seq);

	INSERT INTO history_seq(num) 
		SELECT 1
		WHERE NOT EXISTS (SELECT 1 from history_seq);

	------------------------------triggers---------------------------		

	CREATE TRIGGER IF NOT EXISTS contact_after_insert
		AFTER INSERT
		ON contact
		FOR EACH ROW
	BEGIN
		UPDATE timeline_seq SET num = (num + 1);
		UPDATE history_seq SET num = (num + 1);
		UPDATE contact
		SET timeline_id = (SELECT num FROM timeline_seq),
			history_id  = (SELECT num FROM history_seq),
			last_stmt   = 0,
			timestamp   = CURRENT_TIMESTAMP
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
		UPDATE timeline_seq SET num = (num + 1);
		UPDATE history_seq SET num = (num + 1);
		UPDATE contact
		SET timeline_id = (SELECT num FROM timeline_seq),
			history_id  = (SELECT num FROM history_seq),
			last_stmt   = 1,
			timestamp   = CURRENT_TIMESTAMP
		WHERE id = old.id;
	END;
	
	-- Mark for delete
	CREATE TRIGGER IF NOT EXISTS contact_before_update_delete
		BEFORE UPDATE OF
			last_stmt
		ON contact
		FOR EACH ROW
	BEGIN
		SELECT RAISE(ABORT, 'Update "last_stmt" not allowed')
		WHERE (new.last_stmt < 0 OR new.last_stmt > 2)
		   OR (old.last_stmt = 2 AND new.last_stmt <> old.last_stmt);
	END;
	
	CREATE TRIGGER IF NOT EXISTS contact_after_update_delete
		AFTER UPDATE OF
			last_stmt
		ON contact
		FOR EACH ROW
		WHEN new.last_stmt = 2
	BEGIN
		UPDATE history_seq SET num = (num + 1);
		UPDATE contact
		SET history_id = (SELECT num FROM history_seq),
			last_stmt  = new.last_stmt,
			timestamp  = CURRENT_TIMESTAMP
		WHERE id = old.id;
	END;
	
	`)
	if err != nil {
		log.Fatal(err)
	}
}
