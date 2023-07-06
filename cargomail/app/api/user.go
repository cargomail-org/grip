package api

import (
	"cargomail/app/repository"
)

type UserApi struct {
	user  repository.UserRepository
}
