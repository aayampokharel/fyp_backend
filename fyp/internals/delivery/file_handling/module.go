package filehandling

import (
	"project/internals/domain/repository"
	"project/internals/domain/service"
	"project/internals/usecase"
)

type Module struct {
	currentMappedTCPPort int
	Controller           *Controller
	ParseFileUseCase     *usecase.ParseFileUseCase
	BlockChainUseCase    *usecase.BlockChainUseCase
}

func NewModule(service service.Service, BlockChainRepo repository.IBlockChainRepository,
	NodeRepo repository.INodeRepository, currentMappedTCPPort int, countPrepareMap, countCommitMap map[int]int, operationCounter *int,
	SqlRepo repository.ISqlRepository, pbftService service.PBFTService) *Module {
	pbftUseCase := usecase.NewPBFTUseCase(service, SqlRepo, NodeRepo, countPrepareMap, countCommitMap, operationCounter, pbftService)
	parseFileUseCase := usecase.NewParseFileUseCase(service, SqlRepo)
	blockChainUseCase := usecase.NewBlockChainUseCase(BlockChainRepo, NodeRepo, SqlRepo, service)
	return &Module{Controller: NewController(parseFileUseCase, blockChainUseCase, currentMappedTCPPort, pbftUseCase), ParseFileUseCase: parseFileUseCase, BlockChainUseCase: blockChainUseCase, currentMappedTCPPort: currentMappedTCPPort}
}
