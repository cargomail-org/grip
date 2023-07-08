package transfer

import (
	"net/http"
)

func (svc *service) routes(mux *http.ServeMux) {
	mux.Handle("/api/v1/health", svc.apis.Health.Healthcheck())
}
