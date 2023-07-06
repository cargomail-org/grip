package api

import (
	"cargomail/app/api/helper"
	"cargomail/app/repository"
	"context"
	"errors"
	"net/http"
	"strings"
)

type Apis struct {
	Health HealthApi
	Form   FormApi
	File   FileApi
	Token  TokenApi
	User   UserApi
}

func NewApis(repository repository.Repository, storagePath, domainName string) Apis {
	return Apis{
		Health: HealthApi{domainName: domainName},
		Form:   FormApi{domainName: domainName},
		File:   FileApi{file: repository.File, storagePath: storagePath},
		Token:  TokenApi{user: repository.User, token: repository.Token},
		User:   UserApi{user: repository.User, token: repository.Token},
	}
}

type contextKey string

const userContextKey = contextKey("user")

func (apis *Apis) contextSetUser(r *http.Request, user *repository.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// middleware
func (apis *Apis) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			helper.ReturnErr(w, repository.ErrInvalidOrMissingAuthToken, http.StatusForbidden)
			return
		}

		token := headerParts[1]

		// TODO magic number!
		if len(token) != 52 {
			helper.ReturnErr(w, repository.ErrInvalidOrMissingAuthToken, http.StatusForbidden)
			return
		}

		user, err := apis.User.user.GetByToken(repository.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrUsernameNotFound):
				helper.ReturnErr(w, repository.ErrInvalidCredentials, http.StatusForbidden)
			default:
				helper.ReturnErr(w, err, http.StatusInternalServerError)
			}
			return
		}

		r = apis.contextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}
