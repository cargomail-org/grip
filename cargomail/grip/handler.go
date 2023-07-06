package grip

import "net/http"

func (svc *Service) Health() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Okay"))
	})
}
