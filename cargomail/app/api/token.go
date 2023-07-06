package api

import (
	"cargomail/app/api/helper"
	"cargomail/app/repository"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type TokenApi struct {
	user  repository.UserRepository
	token repository.TokenRepository
}

func (api *TokenApi) Authenticate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var input credentials

		helper.FromJson(r.Body, &input)

		if !validCredentials(input) {
			helper.ReturnErr(w, repository.ErrInvalidCredentialsFormat, http.StatusInternalServerError)
			return
		}

		user, err := api.user.GetByUsername(input.Username)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrUsernameNotFound):
				helper.ReturnErr(w, err, http.StatusForbidden)
			default:
				helper.ReturnErr(w, err, http.StatusInternalServerError)
			}
			return
		}

		match, err := user.Password.Matches(input.Password)
		if err != nil {
			helper.ReturnErr(w, err, http.StatusInternalServerError)
			return
		}

		if !match {
			helper.ReturnErr(w, repository.ErrInvalidCredentials, http.StatusForbidden)
			return
		}

		token, err := api.token.New(user.ID, 24*time.Hour, repository.ScopeAuthentication)
		if err != nil {
			helper.ReturnErr(w, err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(token)
	})
}

func (api *TokenApi) Logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			helper.ReturnErr(w, repository.ErrInvalidOrMissingAuthToken, http.StatusForbidden)
			return
		}

		token := headerParts[1]

		err := api.token.Remove(token)
		if err != nil {
			helper.ReturnErr(w, err, http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
