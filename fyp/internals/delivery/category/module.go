package category

import (
	"project/internals/domain/repository"
	"project/internals/domain/service"
	"project/internals/usecase"
)

type Module struct {
	Controller *Controller
	sqlUseCase *usecase.SqlUseCase
}

func NewModule(sqlRepo repository.ISqlRepository, service service.Service) *Module {
	sqlUseCase := usecase.NewSqlUseCase(sqlRepo, service)
	return &Module{
		Controller: NewController(sqlUseCase),
		sqlUseCase: sqlUseCase,
	}
}
