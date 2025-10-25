package sse

import (
	"project/internals/domain/repository"
	"project/internals/domain/service"
	"project/internals/usecase"
)

type Module struct {
	Controller *Controller
	UseCase    *usecase.SqlUseCase
}

func NewModule(sqlRepo repository.ISqlRepository) *Module {
	service := service.Service{}
	uc := usecase.NewSqlUseCase(sqlRepo, service)
	controller := NewController(*uc)

	return &Module{
		Controller: controller,
		UseCase:    uc,
	}
}
