package provider

import (
	"net/http"
)

func (svc *service) routes(mux *http.ServeMux) {
	// App
	mux.Handle("/", svc.app.Authenticate(svc.app.HomePage()))
	mux.Handle("/login", svc.app.LoginPage())
	mux.Handle("/logout", svc.app.Logout())
	mux.Handle("/register", svc.app.RegisterPage())

	// mux.Handle("/auth/authenticate", svc.app.Session.Authenticate())
	// Health API
	mux.Handle("/api/v1/health", svc.api.Health.Healthcheck())
	// Auth API
	mux.Handle("/api/v1/auth/register", svc.api.User.Register())
	mux.Handle("/api/v1/auth/authenticate", svc.api.Session.Login())
	mux.Handle("/api/v1/auth/logout", svc.api.Session.Logout())
	// Resources API
	mux.Handle("/api/v1/resources/upload", svc.api.Authenticate(svc.api.Resources.Upload()))
	mux.Handle("/api/v1/resources", svc.api.Authenticate(svc.api.Resources.GetAll()))
}
