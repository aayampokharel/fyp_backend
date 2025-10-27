package sse

import (
	"project/internals/domain/repository"
	"project/internals/domain/service"
	"project/internals/usecase"
)

type Module struct {
	Controller *Controller
	SqlUseCase *usecase.SqlUseCase
	SseUseCase *usecase.SSEUseCase
}

func NewModule(sqlRepo repository.ISqlRepository, sseManager *service.SSEManager, sseUseCase *usecase.SSEUseCase) *Module {
	service := *service.NewService()
	uc := usecase.NewSqlUseCase(sqlRepo, service)
	ssuc := usecase.NewSSEUseCase(sqlRepo, sseManager)

	return &Module{
		Controller: NewController(uc, ssuc),
		SqlUseCase: uc,
		SseUseCase: sseUseCase,
	}
}
