package authentication

import (
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	"project/internals/domain/service"
	"project/internals/usecase"
)

type Module struct {
	Controller *Controller
	UseCase    *usecase.SqlUseCase
	FacultyCh  chan<- entity.Institution
	ChannelMap map[string]chan<- entity.Institution
	SSEService *service.SSEManager
}

func NewModule(sqlRepo repository.ISqlRepository, facultyCh chan<- entity.Institution, channelMap map[string]chan<- entity.Institution, sseService *service.SSEManager) *Module {
	service := service.Service{}
	uc := usecase.NewSqlUseCase(sqlRepo, service)
	controller := NewController(*uc)

	return &Module{
		Controller: controller,
		UseCase:    uc,
		FacultyCh:  facultyCh,
		ChannelMap: channelMap,
		SSEService: sseService,
	}
}
