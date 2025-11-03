package delivery

import (
	"project/internals/domain/repository"
	"project/internals/domain/service"
	"project/internals/usecase"
)

type Module struct {
	Controller *Controller
	UseCase    *usecase.BlockChainUseCase
}

func NewModule(blockRepo repository.IBlockChainRepository, nodeRepo repository.INodeRepository, sqlRepo repository.ISqlRepository) *Module {
	blockchainService := *service.NewService()
	uc := usecase.NewBlockChainUseCase(blockRepo, nodeRepo, sqlRepo, blockchainService)
	sqlUc := usecase.NewSqlUseCase(sqlRepo, blockchainService)
	parseUc := usecase.NewParseFileUseCase(blockchainService)
	controller := NewController(*uc, parseUc, sqlUc)

	return &Module{
		Controller: controller,
		UseCase:    uc,
	}
}
