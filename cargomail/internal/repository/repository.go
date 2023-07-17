package repository

import (
	"database/sql"
	"errors"

	tus "github.com/tus/tusd/v2/pkg/handler"
)

type Repository struct {
	Files   FilesRepository
	Storage StorageRepository
	Session SessionRepository
	User    UserRepository
}

var (
	ErrUsernameAlreadyTaken     = errors.New("username already taken")
	ErrUsernameNotFound         = errors.New("username not found")
	ErrInvalidCredentials       = errors.New("invalid authentication credentials")
	ErrMissingUserContext       = errors.New("missing user context")
	ErrInvalidOrMissingSession  = errors.New("invalid or missing session")
	ErrFailedValidationResponse = errors.New("failed validation")
	ErrFileNameNotFound         = errors.New("filename not found")
)

func NewRepository(db *sql.DB, tusHandler *tus.Handler) Repository {
	return Repository{
		Files:   FilesRepository{db: db},
		Storage: StorageRepository{db: db, tusHandler: tusHandler},
		Session: SessionRepository{db: db},
		User:    UserRepository{db: db},
	}
}
