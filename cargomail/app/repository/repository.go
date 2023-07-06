package repository

import (
	"database/sql"
	"errors"
)

type Repository struct {
	File  FileRepository
	Token TokenRepository
	User  UserRepository
}

var (
	ErrUsernameAlreadyTaken      = errors.New("username already taken")
	ErrUsernameNotFound          = errors.New("username not found")
	ErrInvalidCredentials        = errors.New("invalid authentication credentials")
	ErrInvalidCredentialsFormat  = errors.New("invalid authentication credentials format")
	ErrMissingUserContext        = errors.New("missing user context")
	ErrInvalidOrMissingAuthToken = errors.New("invalid or missing authentication token")
	ErrFailedValidationResponse  = errors.New("failed validation")
)

func NewRepository(db *sql.DB) Repository {
	return Repository{
		File:  FileRepository{db: db},
		Token: TokenRepository{db: db},
		User:  UserRepository{db: db},
	}
}
