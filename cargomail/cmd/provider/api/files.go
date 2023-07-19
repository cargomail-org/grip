package api

import (
	"cargomail/cmd/provider/api/helper"
	"cargomail/internal/repository"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/google/uuid"
)

type FilesApi struct {
	files     repository.FilesRepository
	filesPath string
}

func (api *FilesApi) Upload() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			user, ok := r.Context().Value(repository.UserContextKey).(*repository.User)
			if !ok {
				helper.ReturnErr(w, repository.ErrMissingUserContext, http.StatusInternalServerError)
				return
			}

			err := r.ParseMultipartForm(32 << 20)
			if err != nil {
				log.Println(err)
				return
			}

			helper.SetJsonHeader(w)
			w.WriteHeader(http.StatusCreated)

			files := r.MultipartForm.File["files"]
			for i := range files {
				file, err := files[i].Open()
				if err != nil {
					fmt.Println(err)
					return
				}
				defer file.Close()

				filepath := api.filesPath

				if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
					err := os.MkdirAll(filepath, os.ModePerm)
					if err != nil {
						log.Println(err)
						return
					}
				}

				uuid := uuid.NewString()

				f, err := os.OpenFile(filepath+uuid, os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer f.Close()

				hash := sha256.New()
				written, err := io.Copy(f, io.TeeReader(file, hash))
				if err != nil {
					log.Println(err)
					return
				}

				checksum := hash.Sum(nil)
				contentType := files[i].Header.Get("content-type")
				createdAt, err := api.files.Create(user, uuid, checksum, files[i].Filename, filepath, contentType, written)
				if err != nil {
					log.Println(err)
					return
				}

				checksumString := fmt.Sprintf("%x", checksum)
				uploadedFile := repository.File{
					UUID:        uuid,
					Hash:        checksumString,
					Name:        files[i].Filename,
					Size:        written,
					ContentType: contentType,
					CreatedAt:   createdAt,
					Timestamp:   createdAt.UnixMilli(),
				}

				json.NewEncoder(w).Encode(uploadedFile)
			}
		}
	})
}

func (api *FilesApi) Download() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(repository.UserContextKey).(*repository.User)
		if !ok {
			helper.ReturnErr(w, repository.ErrMissingUserContext, http.StatusInternalServerError)
			return
		}

		uuid := path.Base(r.URL.Path)

		fileName, err := api.files.GetOriginalFileName(user, uuid)
		if err != nil {
			helper.ReturnErr(w, err, http.StatusNotFound)
			return
		}

		if r.Method == "HEAD" {
			w.WriteHeader(http.StatusOK)
		} else if r.Method == "GET" {
			filePath := api.filesPath + uuid
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
			http.ServeFile(w, r, filePath)
		}
	})
}

// TODO validate sort fields/filters.SortSafelist
func validFilters(f repository.Filters) bool {
	return f.Page > 0 ||
		f.Page <= 10_000_000 ||
		f.PageSize > 0 ||
		f.PageSize <= 100
}

func (api *FilesApi) GetAll() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			user, ok := r.Context().Value(repository.UserContextKey).(*repository.User)
			if !ok {
				helper.ReturnErr(w, repository.ErrMissingUserContext, http.StatusInternalServerError)
				return
			}

			qs := r.URL.Query()

			filters := repository.Filters{}

			var err error
			filters.Page, err = strconv.Atoi(qs.Get("page"))
			if err != nil {
				filters.Page = 0
			}

			filters.PageSize, err = strconv.Atoi(qs.Get("page_size"))
			if err != nil {
				filters.PageSize = 0
			}

			filters.Sort = qs.Get("sort")
			if len(filters.Sort) == 0 {
				filters.Sort = "id"
			}

			filters.SortSafelist = []string{"id", "name", "size", "content_type", "created_at", "-id", "-name", "-size", "-content_type", "-created_at"}

			if !validFilters(filters) {
				helper.ReturnErr(w, repository.ErrFailedValidationResponse, http.StatusUnprocessableEntity)
				return
			}

			files, _, err := api.files.GetAll(user, filters)
			if err != nil {
				helper.ReturnErr(w, err, http.StatusInternalServerError)
				return
			}

			// log.Printf("metadata: %v", metadata)

			helper.SetJsonHeader(w)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(files)
		}
	})
}

func (api *FilesApi) DeleteByUuidList() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			user, ok := r.Context().Value(repository.UserContextKey).(*repository.User)
			if !ok {
				helper.ReturnErr(w, repository.ErrMissingUserContext, http.StatusInternalServerError)
				return
			}

			var l []string

			err := json.NewDecoder(r.Body).Decode(&l)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = api.files.DeleteByUuidList(user, l)
			if err != nil {
				helper.ReturnErr(w, err, http.StatusInternalServerError)
				return
			}

			filepath := api.filesPath

			for _, uuid := range l {
				_ = os.Remove(filepath + uuid)
			}

			helper.SetJsonHeader(w)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
		}
	})
}
