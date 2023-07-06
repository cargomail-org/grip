package api

import (
	"cargomail/app/api/helper"
	"cargomail/app/repository"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
)

type FileApi struct {
	file        repository.FileRepository
	storagePath string
}

func (api *FileApi) Upload() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(userContextKey).(*repository.User)
		if !ok {
			helper.ReturnErr(w, repository.ErrMissingUserContext, http.StatusInternalServerError)
			return
		}

		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)

		filepath := api.storagePath

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

		api.file.Create(user, uuid, handler.Filename, filepath, handler.Header.Get("content-type"), written)
	})
}

// TODO validate sort fields/filters.SortSafelist
func validFilters(f repository.Filters) bool {
	return f.Page > 0 ||
		f.Page <= 10_000_000 ||
		f.PageSize > 0 ||
		f.PageSize <= 100
}

func (api *FileApi) List() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(userContextKey).(*repository.User)
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

		files, metadata, err := api.file.GetAll(user, filters)
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
