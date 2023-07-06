package api

import (
	"html/template"
	"net/http"
)

type FormApi struct {
	domainName string
}

func (api *FormApi) HomeForm(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			type Devicevalue_view struct {
				Name string
				Size string
				Type string
			}

			type Data struct {
				Files []Devicevalue_view
			}

			data := Data{}
			for i := 1; i < 10; i++ {
				view := Devicevalue_view{
					Name: "devicetype",
					Size: "iddevice",
					Type: "devicename",
				}

				data.Files = append(data.Files, view)
			}

			t.Execute(w, data)
		}
	})
}

func (api *FormApi) LoginForm(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			view := struct {
				DomainName string
			}{
				api.domainName,
			}

			t.Execute(w, view)
		}
	})
}

func (api *FormApi) RegisterForm(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			view := struct {
				DomainName string
			}{
				api.domainName,
			}

			t.Execute(w, view)
		}
	})
}
