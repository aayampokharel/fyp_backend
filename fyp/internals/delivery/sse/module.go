package sse

import (
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	"project/internals/domain/service"
	"project/internals/usecase"
)

type Module struct {
	Controller    *Controller
	UseCase       *usecase.SqlUseCase
	InstitutionCh <-chan entity.Institution
}

func NewModule(sqlRepo repository.ISqlRepository, insstitutionCh <-chan entity.Institution) *Module {
	service := service.Service{}
	uc := usecase.NewSqlUseCase(sqlRepo, service)

	return &Module{
		Controller:    NewController(uc),
		UseCase:       uc,
		InstitutionCh: insstitutionCh,
	}
}
