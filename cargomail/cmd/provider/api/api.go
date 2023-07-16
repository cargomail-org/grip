package api

import (
	"cargomail/cmd/provider/api/helper"
	"cargomail/internal/repository"
	"context"
	"errors"
	"log"
	"net/http"

	tus "github.com/tus/tusd/v2/pkg/handler"
)

type ApiParams struct {
	DomainName string
	FilesPath  string
	Repository repository.Repository
	TusHandler *tus.Handler
}

type Api struct {
	Health  HealthApi
	Files   FilesApi
	Session SessionApi
	User    UserApi
}

func NewApi(params ApiParams) Api {
	return Api{
		Health:  HealthApi{domainName: params.DomainName},
		Files:   FilesApi{files: params.Repository.Files, filesPath: params.FilesPath, tusHandler: params.TusHandler},
		Session: SessionApi{user: params.Repository.User, session: params.Repository.Session},
		User:    UserApi{user: params.Repository.User},
	}
}

// type contextKey string

// const UserContextKey = contextKey("user")

func (api *Api) contextSetUser(r *http.Request, user *repository.User) *http.Request {
	ctx := context.WithValue(r.Context(), repository.UserContextKey, user)
	return r.WithContext(ctx)
}

// middleware
func (api *Api) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Add("Vary", "Authorization")

		// authorizationHeader := r.Header.Get("Authorization")

		// headerParts := strings.Split(authorizationHeader, " ")
		// if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		// 	helper.ReturnErr(w, repository.ErrInvalidOrMissingAuthToken, http.StatusForbidden)
		// 	return
		// }

		// token := headerParts[1]

		cookie, err := r.Cookie("session")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				http.Error(w, "cookie not found", http.StatusBadRequest)
			default:
				log.Println(err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
			return
		}

		session := cookie.Value

		// TODO magic number!
		if len(session) != 52 {
			helper.ReturnErr(w, repository.ErrInvalidOrMissingSession, http.StatusForbidden)
			return
		}

		user, err := api.User.user.GetBySession(repository.ScopeAuthentication, session)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrUsernameNotFound):
				helper.ReturnErr(w, repository.ErrInvalidCredentials, http.StatusForbidden)
			default:
				helper.ReturnErr(w, err, http.StatusInternalServerError)
			}
			return
		}

		r = api.contextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}
