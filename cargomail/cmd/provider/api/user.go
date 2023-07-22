package api

import (
	"cargomail/internal/repository"
	"net/http"
)

type UserApi struct {
	user repository.UserRepository
}

func (api *UserApi) Profile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log.Println("profile request")

		w.WriteHeader(http.StatusOK)
	})
}
