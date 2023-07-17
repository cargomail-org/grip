package repository

import (
	"database/sql"
	"log"

	tus "github.com/tus/tusd/v2/pkg/handler"
)

type StorageRepository struct {
	db         *sql.DB
	tusHandler *tus.Handler
}

type contextKey string

const UserContextKey = contextKey("user")

func (r StorageRepository) TusServe() {
	go func() {
		for {
			event := <-r.tusHandler.CompleteUploads
			ctx := event.Context
			user, ok := ctx.Value(UserContextKey).(*User)
			if !ok {
				log.Println("tus context error")
			}

			log.Println(user.Username)
		}
	}()
}

