package grip

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

type Service struct {
	DB           *sql.DB
	StoragePath  string
	DomainName   string
	GripApiBind  string
	GripCertFile string
	GripKeyFile  string
}

func (svc *Service) Serve(ctx context.Context, errs *errgroup.Group) {
	mux := http.NewServeMux()
	mux.Handle("/health", svc.Health())

	http1Server := &http.Server{Handler: mux, Addr: svc.GripApiBind}

	errs.Go(func() error {
		<-ctx.Done()
		gracefulStop, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelShutdown()

		err := http1Server.Shutdown(gracefulStop)
		if err != nil {
			return err
		}
		log.Print("GRIP service shutdown gracefully")
		return nil
	})

	errs.Go(func() error {
		log.Printf("GRIP service is listening on https://%s", http1Server.Addr)
		return http1Server.ListenAndServeTLS(svc.GripCertFile, svc.GripKeyFile)
	})
}
