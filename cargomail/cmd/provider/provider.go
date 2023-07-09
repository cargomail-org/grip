package provider

import (
	"cargomail/cmd/provider/api"
	"cargomail/cmd/provider/app"
	"cargomail/internal/repository"
	"context"
	"database/sql"
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

type ServiceParams struct {
	DomainName    string
	ResourcesPath string
	DB            *sql.DB
	ProviderBind  string
}

type service struct {
	app          app.App
	api          api.Api
	providerBind string
}

func NewService(params *ServiceParams) (service, error) {
	repository := repository.NewRepository(params.DB)

	templates, err := LoadTemplates()
	if err != nil {
		return service{}, err
	}

	return service{
		app: app.NewApp(
			app.AppParams{
				DomainName: params.DomainName,
				Repository: repository,
				Templates:  templates,
			}),
		api: api.NewApi(
			api.ApiParams{
				DomainName:    params.DomainName,
				Repository:    repository,
				ResourcesPath: params.ResourcesPath,
			}),
		providerBind: params.ProviderBind,
	}, err
}

const (
	publicDir    = "public"
	templatesDir = "templates"
	layoutsDir   = "templates/layouts"
	extension    = "/*.html"

	HomePage     = "home.page.html"
	LoginPage    = "login.page.html"
	RegisterPage = "register.page.html"
)

var (
	//go:embed public/* templates/* templates/layouts/*
	files embed.FS
	// templates map[string]*template.Template
)

func LoadTemplates() (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)

	tmplFiles, err := fs.ReadDir(files, templatesDir)
	if err != nil {
		return nil, err
	}

	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}

		pt, err := template.ParseFS(files, templatesDir+"/"+tmpl.Name(), layoutsDir+extension)
		if err != nil {
			return nil, err
		}

		templates[tmpl.Name()] = pt
	}

	_, exists := templates[HomePage]
	if !exists {
		log.Printf("template %s not found", HomePage)
		return nil, err
	}

	_, exists = templates[LoginPage]
	if !exists {
		log.Printf("template %s not found", LoginPage)
		return nil, err
	}

	_, exists = templates[RegisterPage]
	if !exists {
		log.Printf("template %s not found", RegisterPage)
		return nil, err
	}

	return templates, nil
}

func (svc *service) Serve(ctx context.Context, errs *errgroup.Group) {
	// Routes
	mux := http.NewServeMux()
	svc.routes(mux, svc.app.Templates)

	fs := http.FileServer(http.FS(files))
	mux.Handle("/"+publicDir+"/", http.StripPrefix("/", fs))

	http1Server := &http.Server{Handler: mux, Addr: svc.providerBind}

	errs.Go(func() error {
		<-ctx.Done()
		gracefulStop, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelShutdown()

		err := http1Server.Shutdown(gracefulStop)
		if err != nil {
			return err
		}
		log.Print("provider service shutdown gracefully")
		return nil
	})

	errs.Go(func() error {
		log.Printf("provider service is listening on http://%s", http1Server.Addr)
		return http1Server.ListenAndServe()
	})
}
