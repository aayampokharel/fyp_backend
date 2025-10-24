package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project/internals/usecase"
	err "project/package/errors"
	"project/package/utils/common"
	"time"

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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
	fmt.Print("hello hello")
	//! to be done in middleware .
	if er := json.NewDecoder(r.Body).Decode(&institution); er != nil {
		c.useCase.Logger.Errorln("[authentication_controller] Error: HandleCreateNewInstitution::", er)
		common.HandleErrorResponse(500, err.ErrDecodingJSONString, er, w)
		return
	}

	if institution.InstitutionName == "" || institution.ToleAddress == "" || institution.DistrictAddress == "" || institution.WardNumber == 0 {
		c.useCase.Logger.Errorln("[authentication_controller] Error: HandleCreateNewInstitution::", errorz.ErrEmptyInstitutionInfo)
		common.HandleErrorResponse(500, err.ErrCreatingInstitutionString, errorz.ErrEmptyInstitutionInfo, w)
		return
	}
	institutionEntity := institution.ToEntity()

	if er := c.useCase.CheckDuplicationByInstitutionInfoUseCase(institutionEntity); er != nil {
		c.useCase.Logger.Errorln("[authentication_controller] Error: HandleCreateNewInstitution::", er)
		common.HandleErrorResponse(500, err.ErrDuplicateString, er, w)
		return

	}

	institutionId, er := c.useCase.InsertInstitutionsUseCase(institutionEntity)
	if er != nil {
		c.useCase.Logger.Errorln("[authentication_controller] Error: HandleCreateNewInstitution::", er)
		common.HandleErrorResponse(500, err.ErrCreatingInstitutionString, er, w)
		return
	}
	common.HandleSuccessResponse(CreateInstutionResponse{
		InstitutionID: institutionId,
		IsActive:      false,
	}, w)

}

func (c *Controller) HandleCreateNewUserAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
	var newUserAccount CreateUserAccountRequest
	if er := json.NewDecoder(r.Body).Decode(&newUserAccount); er != nil {
		common.HandleErrorResponse(500, err.ErrDecodingJSONString, er, w)
	}

	newUserAccountEntity := newUserAccount.ToEntity()

	userAccountID, createdAtStr, er := c.useCase.InsertUserAccountUseCase(newUserAccountEntity, newUserAccount.InstitutionID, newUserAccount.InstitutionLogoBase64)

	if er != nil {
		common.HandleErrorResponse(500, err.ErrCreatingUserAccountString, er, w)
		return
	}
	createdAt, er := time.Parse("2006-01-02 15:04:05", createdAtStr)
	if er != nil {
		common.HandleErrorResponse(500, err.ErrCreatingUserAccountString, er, w)
		return
	}
	common.HandleSuccessResponse(CreateUserAccountResponse{
		UserAccountID: userAccountID,
		CreatedAt:     createdAt.Format(time.RFC3339),
	}, w)

}
