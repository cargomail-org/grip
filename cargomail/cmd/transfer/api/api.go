package api

import "cargomail/internal/repository"

type ApiParams struct {
	DomainName    string
	ResourcesPath string
	Repository    repository.Repository
}

type Api struct {
	Health HealthApi
}

func NewApi(params ApiParams) Api {
	return Api{
		Health: HealthApi{params.DomainName},
	}
}
