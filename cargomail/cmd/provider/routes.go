package provider

import (
	"html/template"
	"net/http"
)

func (svc *service) routes(mux *http.ServeMux, templates map[string]*template.Template) {
	// App
	mux.Handle("/", svc.api.Form.HomeForm(templates[HomePage]))
	mux.Handle("/login", svc.api.Form.LoginForm(templates[LoginPage]))
	mux.Handle("/logout", http.RedirectHandler("/login", http.StatusSeeOther))
	mux.Handle("/register", svc.api.Form.RegisterForm(templates[RegisterPage]))
	// Health API
	mux.Handle("/api/v1/health", svc.api.Health.Healthcheck())
	// Auth API
	mux.Handle("/api/v1/auth/register", svc.api.User.Register())
	mux.Handle("/api/v1/auth/authenticate", svc.api.Session.Authenticate())
	mux.Handle("/api/v1/auth/logout", svc.api.Session.Logout())
	// Resources API
	mux.Handle("/api/v1/resources/upload", svc.api.Authenticate(svc.api.Resources.Upload()))
	mux.Handle("/api/v1/resources", svc.api.Authenticate(svc.api.Resources.GetAll()))
}
