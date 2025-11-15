package filehandling

import (
	"project/constants"
	"project/internals/data/config"
	"project/internals/domain/entity"
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
	SqlRepo repository.ISqlRepository, pbftService service.PBFTService, operationChannelMap map[int]chan entity.PBFTExecutionResultEntity, env *config.Env) *Module {

	pbftUseCase := usecase.NewPBFTUseCase(service, SqlRepo, NodeRepo, countPrepareMap, countCommitMap, operationCounter, pbftService, BlockChainRepo, operationChannelMap)
	parseFileUseCase := usecase.NewParseFileUseCase(service, env, SqlRepo)
	blockChainUseCase := usecase.NewBlockChainUseCase(BlockChainRepo, NodeRepo, SqlRepo, service)
	sqlUseCase := usecase.NewSqlUseCase(SqlRepo, service)
	pingyUrl := env.GetValueForKey(constants.PinggyQrUrl)

	return &Module{Controller: NewController(parseFileUseCase, blockChainUseCase, currentMappedTCPPort, pbftUseCase, pingyUrl, sqlUseCase), ParseFileUseCase: parseFileUseCase, BlockChainUseCase: blockChainUseCase, currentMappedTCPPort: currentMappedTCPPort}
}
