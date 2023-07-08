package provider

import (
	"html/template"
	"net/http"
)

func (svc *service) routes(mux *http.ServeMux, templates map[string]*template.Template) {
	//
	mux.Handle("/", svc.apis.Form.HomeForm(templates[HomePage]))
	mux.Handle("/login", svc.apis.Form.LoginForm(templates[LoginPage]))
	mux.Handle("/logout", http.RedirectHandler("/login", http.StatusSeeOther))
	mux.Handle("/register", svc.apis.Form.RegisterForm(templates[RegisterPage]))
	//
	mux.Handle("/api/v1/health", svc.apis.Health.Healthcheck())
	//
	mux.Handle("/api/v1/auth/register", svc.apis.User.Register())
	mux.Handle("/api/v1/auth/authenticate", svc.apis.Session.Authenticate())
	mux.Handle("/api/v1/auth/logout", svc.apis.Session.Logout())
	//
	mux.Handle("/api/v1/resources/upload", svc.apis.Authenticate(svc.apis.Resources.Upload()))
	mux.Handle("/api/v1/resources", svc.apis.Authenticate(svc.apis.Resources.GetAll()))
}
