package delivery

import (
	"project/internals/data/config"
	"project/internals/domain/repository"
	"project/internals/domain/service"
	"project/internals/usecase"
)

type Module struct {
	Controller *Controller
	UseCase    *usecase.BlockChainUseCase
}

func NewModule(blockRepo repository.IBlockChainRepository, nodeRepo repository.INodeRepository, sqlRepo repository.ISqlRepository, env *config.Env) *Module {
	service := *service.NewService()
	uc := usecase.NewBlockChainUseCase(blockRepo, nodeRepo, sqlRepo, service)
	sqlUc := usecase.NewSqlUseCase(sqlRepo, service)
	parseUc := usecase.NewParseFileUseCase(service, env, sqlRepo)
	controller := NewController(*uc, parseUc, sqlUc)

	return &Module{
		Controller: controller,
		UseCase:    uc,
	}
}
