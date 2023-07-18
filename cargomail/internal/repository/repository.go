package repository

import (
	"database/sql"
	"errors"
)

type Repository struct {
	Files   FilesRepository
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

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Files:   FilesRepository{db: db},
		Session: SessionRepository{db: db},
		User:    UserRepository{db: db},
	}
}
