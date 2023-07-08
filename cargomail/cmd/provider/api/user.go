package api

import (
	"cargomail/internal/repository"
)

type UserApi struct {
	user repository.UserRepository
}
