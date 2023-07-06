package provider

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sync/errgroup"

	"cargomail/app"
	"cargomail/cdb"
	"cargomail/config"
	"cargomail/grip"
)

func Start() error {
	startFlags := config.NewStartFlags()

	ctx, done := context.WithCancel(context.Background())
	defer done()
	errs, ctx := errgroup.WithContext(ctx)

	log.Printf("using database %v", startFlags.DbPath)
	db, err := sql.Open("sqlite3", startFlags.DbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cdb.Init(db)

	// App
	appService := app.NewService(db, startFlags.StoragePath, startFlags.DomainName, startFlags.AppApiBind)
	appService.Serve(ctx, errs)

	// GRIP
	gripService := grip.Service{
		DomainName:   startFlags.DomainName,
		DB:           db,
		StoragePath:  startFlags.StoragePath,
		GripApiBind:  startFlags.GripApiBind,
		GripCertFile: startFlags.GripCertFile,
		GripKeyFile:  startFlags.GripKeyFile,
	}
	gripService.Serve(ctx, errs)

	go func() error {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-stop:
			log.Printf("Received signal: %s\n", sig)
			done()
		case <-ctx.Done():
			return ctx.Err()
		}
		return nil
	}()

	if err := errs.Wait(); err == nil || err == context.Canceled || err == http.ErrServerClosed {
		return nil
	} else {
		return err
	}
}
