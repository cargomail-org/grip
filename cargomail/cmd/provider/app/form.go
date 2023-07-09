package app

import (
	"html/template"
	"net/http"
)

func (app *App) HomeForm(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			type Data struct {
				Username string
			}

			data := Data{Username: "igor"}

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

func (app *App) LoginForm(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			data := struct {
				DomainName string
			}{
				app.domainName,
			}

			t.Execute(w, data)
		}
	})
}

func (app *App) RegisterForm(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			data := struct {
				DomainName string
			}{
				app.domainName,
			}

			t.Execute(w, data)
		}
	})
}
