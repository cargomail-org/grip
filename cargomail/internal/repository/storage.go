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
			// TODO send errors to the client
			event := <-r.tusHandler.CompleteUploads
			ctx := event.Context
			user, ok := ctx.Value(UserContextKey).(*User)
			if !ok {
				log.Println("tus context error")
			}

			log.Println(user.Username)

			// user, ok := event.Context().Value(userContextKey).(*repository.User)

			// id := event.Upload.ID
			// downloadUrl := config.Mailbox.Uri + config.Filestore.BasePath + id
			// filename := event.Upload.MetaData["filename"]
			// mimeType := event.Upload.MetaData["filetype"]
			// uploadId := event.Upload.MetaData["uploadId"]
			// fileSize := event.Upload.Size
			// path := event.Upload.Storage["Path"]
			// uploadSha256sum := event.HTTPRequest.Header.Get("sha256sum")

			// file := resourcev1.File{DownloadUrl: downloadUrl, Filename: filename, MimeType: mimeType, FileSize: fileSize}

			// log.Printf("User %s uploaded %s file %s using uploadId: %s", username, id, mimeType, uploadId)

			// dbFile, err := emailRepository.Repo.FilesCreate(repo, username, &file)
			// if err != nil {
			// 	log.Printf("Files database create error %s", err.Error())
			// }

			// sha256sum, err := checksum(path)
			// if err != nil {
			// 	log.Printf("Checksum error %s", err.Error())
			// }

			// if len(uploadSha256sum) != 64 {
			// 	log.Printf("Checksum %s is not valid on filename %s", uploadSha256sum, filename)
			// }

			// if sha256sum != uploadSha256sum {
			// 	log.Printf("Checksum mismatch on filename %s: %s vs %s", filename, sha256sum, uploadSha256sum)
			// }

			// log.Printf("Checksum: %s", sha256sum)

			// if dbFile != nil {
			// 	dbFile.Sha256Sum = sha256sum

			// 	dbId, err := strconv.ParseInt(dbFile.Id, 10, 64)
			// 	if err != nil {
			// 		log.Printf("Files database id error %s", err.Error())
			// 	}

			// 	_, err = emailRepository.Repo.FilesUpdate(repo, username, dbId, uploadId, dbFile)
			// 	if err != nil {
			// 		log.Printf("Files database update error %s", err.Error())
			// 	}
			// }
		}
	}()
}
