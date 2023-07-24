package repository

import (
	"database/sql"
)

type ContactsRepository struct {
	db *sql.DB
}
