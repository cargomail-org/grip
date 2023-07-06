package app

import (
	"html/template"
	"net/http"
)

func (svc *service) routes(mux *http.ServeMux, templates map[string]*template.Template) {
	mux.Handle("/", svc.apis.Form.HomeForm(templates[HomePage]))
	mux.Handle("/login", svc.apis.Form.LoginForm(templates[LoginPage]))
	mux.Handle("/logout", http.RedirectHandler("/login", http.StatusSeeOther))
	mux.Handle("/register", svc.apis.Form.RegisterForm(templates[RegisterPage]))
	mux.Handle("/api/v1/files/upload", svc.apis.Authenticate(svc.apis.File.Upload()))
	mux.Handle("/api/v1/files", svc.apis.Authenticate(svc.apis.File.List()))

	mux.Handle("/health", svc.apis.Health.Healthcheck())

	mux.Handle("/api/v1/auth/register", svc.apis.User.Register())
	mux.Handle("/api/v1/auth/authenticate", svc.apis.Token.Authenticate())
	mux.Handle("/api/v1/auth/logout", svc.apis.Token.Logout())
}
