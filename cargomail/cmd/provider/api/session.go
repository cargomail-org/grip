package api

import (
	"cargomail/cmd/provider/api/helper"
	"cargomail/internal/repository"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

type SessionApi struct {
	user    repository.UserRepository
	session repository.SessionRepository
}

type credentials struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
}

func isAlnumOrHyphenOrDot(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '-' && r != '.' {
			return false
		}
	}
	return true
}

func validCredentials(c credentials) bool {
	return len(c.Username) <= 40 &&
		len(c.Username) > 0 &&
		!strings.HasPrefix(c.Username, "-") &&
		!strings.Contains(c.Username, "--") &&
		!strings.HasSuffix(c.Username, "-") &&
		!strings.HasPrefix(c.Username, ".") &&
		!strings.Contains(c.Username, "..") &&
		!strings.HasSuffix(c.Username, ".") &&
		isAlnumOrHyphenOrDot(c.Username) &&
		len(c.Password) <= 40 &&
		len(c.Password) > 0
}

func (api *UserApi) Register() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var input credentials

		helper.FromJson(r.Body, &input)

		if !validCredentials(input) {
			helper.ReturnErr(w, repository.ErrInvalidCredentials, http.StatusForbidden)
			return
		}

		user := &repository.User{
			Username: input.Username,
		}

		err := user.Password.Set(input.Password)
		if err != nil {
			helper.ReturnErr(w, err, http.StatusInternalServerError)
			return
		}

		err = api.user.Create(user)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrUsernameAlreadyTaken):
				helper.ReturnErr(w, err, http.StatusForbidden)
			default:
				helper.ReturnErr(w, err, http.StatusInternalServerError)
			}
			return
		}

		helper.SetJsonHeader(w)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	})
}

func (api *SessionApi) Login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var input credentials

		helper.FromJson(r.Body, &input)

		if !validCredentials(input) {
			helper.ReturnErr(w, repository.ErrInvalidCredentials, http.StatusForbidden)
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

		session, err := api.session.New(user.ID, 24*time.Hour, repository.ScopeAuthentication)
		if err != nil {
			helper.ReturnErr(w, err, http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{
			Name:     "session",
			Value:    session.Plaintext,
			Path:     "/",
			HttpOnly: true,
			Secure:   false, // !!!
			SameSite: http.SameSiteLaxMode,
		}

		if input.RememberMe {
			cookie.Expires = session.Expiry
		}

		http.SetCookie(w, &cookie)

		helper.SetJsonHeader(w)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(session)
	})
}

func (api *SessionApi) Logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clearCookie := http.Cookie{
			Name:     "session",
			Value:    "",
			MaxAge:   -1,
			Path:     "/",
			HttpOnly: true,
			Secure:   false, // !!!
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, &clearCookie)

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

		err = api.session.Remove(session)
		if err != nil {
			helper.ReturnErr(w, err, http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
