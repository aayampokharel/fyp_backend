package filehandling

import (
	"project/internals/domain/repository"
	"project/internals/domain/service"
	"project/internals/usecase"
)

type Module struct {
	Controller        *Controller
	ParseFileUseCase  *usecase.ParseFileUseCase
	BlockChainUseCase *usecase.BlockChainUseCase
}

func NewModule(service service.Service, BlockChainRepo repository.IBlockChainRepository,
	NodeRepo repository.INodeRepository,
	SqlRepo repository.ISqlRepository) *Module {

	parseFileUseCase := usecase.NewParseFileUseCase(service, SqlRepo)
	blockChainUseCase := usecase.NewBlockChainUseCase(BlockChainRepo, NodeRepo, SqlRepo, service)
	return &Module{Controller: NewController(parseFileUseCase, blockChainUseCase), ParseFileUseCase: parseFileUseCase, BlockChainUseCase: blockChainUseCase}
}
