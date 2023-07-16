package transfer

import (
	"cargomail/cmd/transfer/api"
	"cargomail/internal/repository"
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/tus/tusd/v2/pkg/filestore"
	tus "github.com/tus/tusd/v2/pkg/handler"
	"golang.org/x/sync/errgroup"
)

type ServiceParams struct {
	DomainName       string
	FilesPath        string
	DB               *sql.DB
	TransferCertPath string
	TransferKeyPath  string
	TransferBind     string
}

type service struct {
	api              api.Api
	transferBind     string
	TransferCertPath string
	TransferKeyPath  string
}

func NewService(params *ServiceParams) service {
	// tus
	store := filestore.FileStore{
		Path: params.FilesPath,
	}

	composer := tus.NewStoreComposer()
	store.UseIn(composer)

	tusHandler, err := tus.NewHandler(tus.Config{
		BasePath:              "/api/v1/files/tus/upload",
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	})
	if err != nil {
		log.Printf("unable to create tus handler: %s", err)
	}

	repository := repository.NewRepository(params.DB, tusHandler)
	return service{
		api: api.NewApi(
			api.ApiParams{
				DomainName: params.DomainName,
				Repository: repository,
				FilesPath:  params.FilesPath,
			}),
		transferBind:     params.TransferBind,
		TransferCertPath: params.TransferCertPath,
		TransferKeyPath:  params.TransferKeyPath,
	}
}

func (svc *service) Serve(ctx context.Context, errs *errgroup.Group) {
	// Routes
	mux := http.NewServeMux()
	svc.routes(mux)

	http1Server := &http.Server{Handler: mux, Addr: svc.transferBind}

	errs.Go(func() error {
		<-ctx.Done()
		gracefulStop, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelShutdown()

		err := http1Server.Shutdown(gracefulStop)
		if err != nil {
			return err
		}
		log.Print("transfer service shutdown gracefully")
		return nil
	})

	errs.Go(func() error {
		log.Printf("transfer service is listening on https://%s", http1Server.Addr)
		return http1Server.ListenAndServeTLS(svc.TransferCertPath, svc.TransferKeyPath)
	})
}
