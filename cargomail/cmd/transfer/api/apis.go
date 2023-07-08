package api

import "cargomail/internal/repository"

type ApisParams struct {
	DomainName    string
	ResourcesPath string
	Repository    repository.Repository
}

type Apis struct {
	Health HealthApi
}

func NewApis(params ApisParams) Apis {
	return Apis{
		Health: HealthApi{params.DomainName},
	}
}
