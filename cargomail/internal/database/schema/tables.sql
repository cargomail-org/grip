	PRAGMA foreign_keys=ON;

	------------------------------tables-----------------------------

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

	CREATE TABLE IF NOT EXISTS contact_timeline_seq (
		user_id 		INTEGER NOT NULL REFERENCES user ON DELETE CASCADE,
		last_timeline_id integer(8) NOT NULL
	);
	CREATE TABLE IF NOT EXISTS contact_history_seq (
		user_id 		INTEGER NOT NULL REFERENCES user ON DELETE CASCADE,
		last_history_id integer(8) NOT NULL
	);

	------------------------------indexes----------------------------

	CREATE INDEX IF NOT EXISTS idx_file_hash ON file(hash);

	CREATE INDEX IF NOT EXISTS idx_contact_timeline_id ON contact (timeline_id);
	CREATE INDEX IF NOT EXISTS idx_contact_history_id ON contact (history_id);
	CREATE INDEX IF NOT EXISTS idx_contact_last_stmt ON contact (last_stmt);

	CREATE UNIQUE INDEX IF NOT EXISTS idx_contact_timeline_seq ON contact_timeline_seq(user_id);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_contact_history_seq ON contact_history_seq(user_id);

