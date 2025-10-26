package admin

import (
	"project/internals/domain/repository"
	"project/internals/domain/service"
	"project/internals/usecase"
)

type Module struct {
	Controller *Controller
	sqlUseCase *usecase.SqlUseCase
	sseUseCase *usecase.SSEUseCase
}

func NewModule(sqlRepo repository.ISqlRepository, service service.Service, sseManager *service.SSEManager) *Module {
	sseUseCase := usecase.NewSSEUseCase(sqlRepo, sseManager)
	sqlUseCase := usecase.NewSqlUseCase(sqlRepo, service)
	return &Module{
		Controller: NewController(sqlUseCase, sseUseCase),
		sqlUseCase: sqlUseCase,
		sseUseCase: sseUseCase,
	}
}
