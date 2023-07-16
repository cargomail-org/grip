package api

import (
	"cargomail/cmd/provider/api/helper"
	"cargomail/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
	tus "github.com/tus/tusd/v2/pkg/handler"
)

type FilesApi struct {
	files      repository.FilesRepository
	filesPath  string
	tusHandler *tus.Handler
}

func (api *FilesApi) TusUpload() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method)
		if r.Method == "POST" {
			api.tusHandler.PostFile(w, r)
		} else if r.Method == "HEAD" {
			api.tusHandler.HeadFile(w, r)
		} else if r.Method == "PATCH" {
			api.tusHandler.PatchFile(w, r)
		} else if r.Method == "DEL" {
			api.tusHandler.DelFile(w, r)
		}
	})
}

func (api *FilesApi) Upload() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusOK)

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

			written, err := io.Copy(f, file)
			if err != nil {
				log.Println(err)
				return
			}

			api.files.Create(user, uuid, files[i].Filename, filepath, files[i].Header.Get("content-type"), written)

			uploadedFile := repository.File{
				Name: files[i].Filename,
			}

			json.NewEncoder(w).Encode(uploadedFile)
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
			filters.Page = 1
		}

		filters.PageSize, err = strconv.Atoi(qs.Get("page_size"))
		if err != nil {
			filters.PageSize = 20
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

		files, metadata, err := api.files.GetAll(user, filters)
		if err != nil {
			helper.ReturnErr(w, err, http.StatusInternalServerError)
			return
		}

		log.Printf("metadata: %v", metadata)

		helper.SetJsonHeader(w)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(files)
	})
}
