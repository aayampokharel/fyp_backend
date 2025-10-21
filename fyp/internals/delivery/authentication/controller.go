package authentication

import (
	"encoding/json"
	"net/http"
	"project/internals/usecase"
	"project/package/utils/common"

	errorz "project/package/errors"
)

type Controller struct {
	useCase usecase.SqlUseCase
}

func NewController(useCase usecase.SqlUseCase) *Controller {
	return &Controller{useCase: useCase}
}

func (c *Controller) HandleCreateNewInstitution(w http.ResponseWriter, r *http.Request) {
	var institution CreateInstitutionRequest

	if er := json.NewDecoder(r.Body).Decode(&institution); er != nil {
		c.useCase.Logger.Errorln("[authentication_controller] Error: HandleCreateNewInstitution::", er)
		common.HandleErrorResponse(500, "error decoding json", er, w)
		return
	}

	if institution.InstitutionName == "" || institution.ToleAddress == "" || institution.DistrictAddress == "" || institution.WardNumber == 0 {
		c.useCase.Logger.Errorln("[authentication_controller] Error: HandleCreateNewInstitution::", errorz.ErrEmptyInstitutionInfo)
		common.HandleErrorResponse(500, "error creating institution", errorz.ErrEmptyInstitutionInfo, w)
		return
	}
	institutionEntity := institution.ToEntity()

	if er := c.useCase.CheckDuplicationByInstitutionInfoUseCase(institutionEntity); er != nil {
		c.useCase.Logger.Errorln("[authentication_controller] Error: HandleCreateNewInstitution::", er)
		common.HandleErrorResponse(500, "duplicate error", er, w)
		return

	}

	institutionId, er := c.useCase.InsertInstitutionsUseCase(institutionEntity)
	if er != nil {
		c.useCase.Logger.Errorln("[authentication_controller] Error: HandleCreateNewInstitution::", er)
		common.HandleErrorResponse(500, "error creating institution", er, w)
		return
	}
	common.HandleSuccessResponse(CreateInstutionResponse{
		InstitutionID: institutionId,
	}, w)

}
