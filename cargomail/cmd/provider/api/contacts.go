package api

import (
	"cargomail/cmd/provider/api/helper"
	"cargomail/internal/repository"
	"encoding/json"
	"log"
	"net/http"
)

type ContactsApi struct {
	contacts repository.ContactsRepository
}

func (api *ContactsApi) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(repository.UserContextKey).(*repository.User)
		if !ok {
			helper.ReturnErr(w, repository.ErrMissingUserContext, http.StatusInternalServerError)
			return
		}

		var contact *repository.Contact

		err := json.NewDecoder(r.Body).Decode(&contact)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		contact, err = api.contacts.Create(user, contact)
		if err != nil {
			log.Println(err)
			return
		}

		helper.SetJsonResponse(w, http.StatusCreated, contact)
	})
}

func (api *ContactsApi) GetAll() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(repository.UserContextKey).(*repository.User)
		if !ok {
			helper.ReturnErr(w, repository.ErrMissingUserContext, http.StatusInternalServerError)
			return
		}

		contacts, err := api.contacts.GetAll(user)
		if err != nil {
			log.Println(err)
			return
		}

		helper.SetJsonResponse(w, http.StatusCreated, contacts)
	})
}

func (api *ContactsApi) Update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func (api *ContactsApi) DeleteByUuidList() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
