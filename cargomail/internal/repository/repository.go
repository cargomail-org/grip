package repository

import (
	"database/sql"
	"errors"
)

type Repository struct {
	Resources ResourcesRepository
	Session   SessionRepository
	User      UserRepository
}

var (
	ErrUsernameAlreadyTaken      = errors.New("username already taken")
	ErrUsernameNotFound          = errors.New("username not found")
	ErrInvalidCredentials        = errors.New("invalid authentication credentials")
	ErrInvalidCredentialsFormat  = errors.New("invalid authentication credentials format")
	ErrMissingUserContext        = errors.New("missing user context")
	ErrInvalidOrMissingSession = errors.New("invalid or missing session")
	ErrFailedValidationResponse  = errors.New("failed validation")
)

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Resources: ResourcesRepository{db: db},
		Session:   SessionRepository{db: db},
		User:      UserRepository{db: db},
	}
}
