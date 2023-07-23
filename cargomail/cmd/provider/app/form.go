package app

import (
	"cargomail/internal/repository"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (app *App) HomePage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			type Data struct {
				DomainName     string
				Username       string
				UsernameLetter string
			}

			user, ok := r.Context().Value(repository.UserContextKey).(*repository.User)
			if !ok {
				log.Fatal("missing user context")
			}

			data := Data{DomainName: app.domainName, Username: user.Username, UsernameLetter: fmt.Sprintf("%c", strings.ToUpper(user.Username)[0])}

			t := app.HomeTemplate
			t.Execute(w, data)
		}
	})
}

// func (app *App) LoginAction(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	username := r.Form.Get("username")
// 	password := r.Form.Get("password")

// 	data := map[string]interface{}{
// 		"err": "invalid credentials",
// 	}

// 	if username == "igor" && password == "password" {
// 		redirectToHomePage(w, r)
// 	} else {
// 		app.Templates["login.page.html"].Execute(w, data)
// 	}
// }

func (app *App) LoginPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			data := struct {
				DomainName string
			}{
				app.domainName,
			}

			t := app.LoginTemplate
			t.Execute(w, data)
		}
	})
}

func (app *App) RegisterPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			data := struct {
				DomainName string
			}{
				app.domainName,
			}

			t := app.RegisterTemplate
			t.Execute(w, data)
		}
	})
}
