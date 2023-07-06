package app

import (
	"cargomail/app/api"
	"cargomail/app/repository"
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

type service struct {
	apis       api.Apis
	appApiBind string
}

func NewService(db *sql.DB, storagePath, domainName, appApiBind string) service {
	repository := repository.NewRepository(db)
	return service{
		apis:       api.NewApis(repository, storagePath, domainName),
		appApiBind: appApiBind,
	}
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
	files     embed.FS
	templates map[string]*template.Template
)

func LoadTemplates() error {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	tmplFiles, err := fs.ReadDir(files, templatesDir)
	if err != nil {
		return err
	}

	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}

		pt, err := template.ParseFS(files, templatesDir+"/"+tmpl.Name(), layoutsDir+extension)
		if err != nil {
			return err
		}

		templates[tmpl.Name()] = pt
	}
	return nil
}

func (svc *service) Serve(ctx context.Context, errs *errgroup.Group) {
	err := LoadTemplates()
	if err != nil {
		log.Fatal(err)
	}

	_, exists := templates[HomePage]
	if !exists {
		log.Printf("template %s not found", HomePage)
		log.Fatal(err)
	}

	_, exists = templates[LoginPage]
	if !exists {
		log.Printf("template %s not found", LoginPage)
		log.Fatal(err)
	}

	_, exists = templates[RegisterPage]
	if !exists {
		log.Printf("template %s not found", RegisterPage)
		log.Fatal(err)
	}

	// Routes
	mux := http.NewServeMux()
	svc.routes(mux, templates)

	fs := http.FileServer(http.FS(files))
	mux.Handle("/"+publicDir+"/", http.StripPrefix("/", fs))

	http1Server := &http.Server{Handler: mux, Addr: svc.appApiBind}

	errs.Go(func() error {
		<-ctx.Done()
		gracefulStop, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelShutdown()

		err := http1Server.Shutdown(gracefulStop)
		if err != nil {
			return err
		}
		log.Print("App service shutdown gracefully")
		return nil
	})

	errs.Go(func() error {
		log.Printf("App service is listening on http://%s", http1Server.Addr)
		return http1Server.ListenAndServe()
	})
}
