package delivery

import (
	"project/internals/domain/repository"
	"project/internals/domain/service"
	"project/internals/usecase"
)

type Module struct {
	Controller *Controller
}

func NewModule(repo repository.IBlockChainRepository) *Module {
	blockchainService := service.Service{}
	usecase := usecase.NewBlockChainUseCase(repo, blockchainService)
	return &Module{Controller: NewController(*usecase)}

}
