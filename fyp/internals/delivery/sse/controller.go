package sse

import (
	"fmt"
	"net/http"
	"project/internals/domain/entity"
	"project/internals/usecase"
	"project/package/enum"
)

type Controller struct {
	useCase usecase.SqlUseCase
}

func NewController(useCase usecase.SqlUseCase) *Controller {
	return &Controller{useCase: useCase}
}

func (c *Controller) SendInstitutionsToBeVerified(newInstitutionCh <-chan entity.Institution, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	institutionsList, er := c.useCase.GetInstitutionsToBeVerifiedUseCase()
	if er != nil {
		return
	}
	HandleSSEResponse(institutionsList, enum.SSEFORMLIST, w)

	for {
		select {
		case newInstitution := <-newInstitutionCh:
			HandleSSEResponse(newInstitution, enum.SSESINGLEFORM, w)
		case <-ctx.Done():
			c.useCase.Logger.Infoln("[send_institutions_to_be_verified] Info: sendInstitutionsToBeVerified::CLIENt disconected ! ", fmt.Sprint(ctx.Err()))
			return
		}
	}
}
