package admin

import (
	"project/internals/domain/entity"
	"project/internals/usecase"
	"project/package/enum"
	err "project/package/errors"
	"project/package/utils/common"
	"time"
)

type Controller struct {
	sqlUseCase *usecase.SqlUseCase
	sseUseCase *usecase.SSEUseCase
}

func NewController(sqlUseCase *usecase.SqlUseCase, sseUseCase *usecase.SSEUseCase) *Controller {
	return &Controller{sqlUseCase: sqlUseCase, sseUseCase: sseUseCase}
}

func (c *Controller) HandleAdminLogin(AdminLoginRequest AdminLoginRequest) entity.Response {
	userID, createdAt, er := c.sqlUseCase.VerifyUserLoginUseCase(AdminLoginRequest.AdminEmail, AdminLoginRequest.Password, enum.ADMIN)

	if er != nil {
		return common.HandleErrorResponse(401, err.ErrVerifyingAdminString, er)
	}
	adminLoginResponse, er := c.sseUseCase.GetAllPendingInstitutionsForAdminsUseCase(userID, createdAt)
	if er != nil {
		return common.HandleErrorResponse(500, err.ErrVerifyingAdminString, er)
	}
	return common.HandleSuccessResponse(AdminLoginResponse{UserID: adminLoginResponse.UserID, CreatedAt: adminLoginResponse.CreatedTime.Format(time.RFC3339), InstitutionList: adminLoginResponse.InstitutionList})

}

func (c *Controller) HandleDeleteSSEClient(token string) entity.Response {
	////used when tab isnot closed but admin logsout .
	if er := c.sseUseCase.RemoveClientUseCase(token); er != nil {
		return common.HandleErrorResponse(500, err.ErrRemovingClientString, er)
	}
	return common.HandleSuccessResponse(nil)
}

func (c *Controller) HandleGetPendingInstitutionList(request map[string]string) entity.Response {
	requestMap, er := common.CheckMapKeysReturnValues(request, GetAllPendingInstitutionsQuery)
	if er != nil {
		return common.HandleErrorResponse(400, err.ErrParsingQueryParametersString, er)
	}
	adminID := requestMap["admin_id"]
	pendingInstitutionList, er := c.sqlUseCase.GetAllPendingInstitutionsUseCase(adminID)
	if er != nil {
		return common.HandleErrorResponse(500, err.ErrGettingPendingInstitutionListString, er)
	}
	return common.HandleSuccessResponse(GetAllPendingInstitutionsResponse{
		PendingInstitutionList: pendingInstitutionList,
	})
}
