// package delivery

// import (
// 	"project/internals/domain/repository"
// 	"project/internals/domain/service"
// 	"project/internals/usecase"
// )

// type Module struct {
// 	Controller *Controller
// }

// func NewModule(repo repository.IBlockChainRepository) *Module {
// 	blockchainService := service.Service{}
// 	usecase := usecase.NewBlockChainUseCase(repo, blockchainService)
// 	return &Module{Controller: NewController(*usecase)}

// }

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
	blockchainService := service.Service{}
	uc := usecase.NewBlockChainUseCase(blockRepo, nodeRepo, sqlRepo, blockchainService)
	controller := NewController(*uc)

	return &Module{
		Controller: controller,
		UseCase:    uc,
	}
}
