package usecase

import (
	"project/internals/domain/repository"
	"project/internals/domain/service"
	"project/package/utils/common"
	logger "project/package/utils/pkg"
	"time"

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

func (uc *SSEUseCase) VerifyAdminLoginUseCase(userEmail, password string) (userID string, generatedUniqueToken string, createdTime time.Time, er error) {
	userID, createdAt, er := uc.SqlRepo.VerifyAdminLogin(userEmail, password)
	if er != nil {
		return "", "", time.Time{}, er
	}
	generatedUniqueToken = common.GenerateUUID(20)
	uc.SSEManager.AddClient(generatedUniqueToken)

	return userID, generatedUniqueToken, createdAt, nil
}

func (uc *SSEUseCase) RemoveClientUseCase(token string) error {
	uc.SSEManager.RemoveClient(token)
	return nil
}
