package admin

import (
	"project/internals/domain/entity"
	"project/internals/usecase"
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
	adminLoginResponse, er := c.sseUseCase.VerifyAdminLoginUseCase(AdminLoginRequest.AdminEmail, AdminLoginRequest.Password)

	if er != nil {
		return common.HandleErrorResponse(500, err.ErrVerifyingAdminString, er)
	}

	return common.HandleSuccessResponse(AdminLoginResponse{UserID: adminLoginResponse.UserID, SSEToken: adminLoginResponse.GeneratedUniqueToken, CreatedAt: adminLoginResponse.CreatedTime.Format(time.RFC3339), InstitutionList: adminLoginResponse.InstitutionList})

}

func (c *Controller) HandleDeleteSSEClient(token string) entity.Response {
	if er := c.sseUseCase.RemoveClientUseCase(token); er != nil {
		return common.HandleErrorResponse(500, err.ErrRemovingClientString, er)
	}
	return common.HandleSuccessResponse(nil)
}
