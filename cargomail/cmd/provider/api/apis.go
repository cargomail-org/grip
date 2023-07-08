package api

import (
	"cargomail/cmd/provider/api/helper"
	"cargomail/internal/repository"
	"context"
	"errors"
	"log"
	"net/http"
)

type ApisParams struct {
	DomainName    string
	ResourcesPath string
	Repository    repository.Repository
}

type Apis struct {
	Health    HealthApi
	Form      FormApi
	Resources ResourcesApi
	Session   SessionApi
	User      UserApi
}

func NewApis(params ApisParams) Apis {
	return Apis{
		Health:    HealthApi{domainName: params.DomainName},
		Form:      FormApi{domainName: params.DomainName},
		Resources: ResourcesApi{resources: params.Repository.Resources, resourcesPath: params.ResourcesPath},
		Session:   SessionApi{user: params.Repository.User, session: params.Repository.Session},
		User:      UserApi{user: params.Repository.User},
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
		// w.Header().Add("Vary", "Authorization")

		// authorizationHeader := r.Header.Get("Authorization")

		// headerParts := strings.Split(authorizationHeader, " ")
		// if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		// 	helper.ReturnErr(w, repository.ErrInvalidOrMissingAuthToken, http.StatusForbidden)
		// 	return
		// }

		// token := headerParts[1]

		cookie, err := r.Cookie("cargomail")
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

		user, err := apis.User.user.GetBySession(repository.ScopeAuthentication, session)
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
