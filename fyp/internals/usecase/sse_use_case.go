package usecase

import (
	"project/internals/domain/repository"
	"project/internals/domain/service"
	"project/internals/usecase/dto"
	logger "project/package/utils/pkg"

	"go.uber.org/zap"
)

type SSEUseCase struct {
	SqlRepo    repository.ISqlRepository
	SSEManager *service.SSEManager
	Logger     *zap.SugaredLogger
}

func NewSSEUseCase(sqlRepo repository.ISqlRepository, sseManager *service.SSEManager) *SSEUseCase {
	return &SSEUseCase{
		SqlRepo:    sqlRepo,
		Logger:     logger.Logger,
		SSEManager: sseManager,
	}
}

func (uc *SSEUseCase) VerifyAdminLoginUseCase(userEmail, password string) (*dto.AdminLoginResponse, error) {
	var adminLoginResponse dto.AdminLoginResponse
	userID, createdAt, er := uc.SqlRepo.VerifyAdminLogin(userEmail, password)
	if er != nil {
		return nil, er
	}
	// generatedUniqueToken := common.GenerateUUID(20)

	// uc.SSEManager.AddClient(generatedUniqueToken)
	institutionList, er := uc.SqlRepo.GetToBeVerifiedInstitutions()
	if er != nil {
		return nil, er
	}
	adminLoginResponse = dto.AdminLoginResponse{UserID: userID, CreatedTime: createdAt, InstitutionList: institutionList}
	return &adminLoginResponse, nil
}

func (uc *SSEUseCase) RemoveClientUseCase(token string) error {
	uc.SSEManager.RemoveClient(token)
	return nil
}
