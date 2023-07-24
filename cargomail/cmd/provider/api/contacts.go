package api

import (
	"cargomail/internal/repository"
	"net/http"
)

type ContactsApi struct {
	contacts repository.ContactsRepository
}

func (api *ContactsApi) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
