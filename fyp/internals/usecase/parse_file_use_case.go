package usecase

import (
	"project/internals/domain/entity"
	"project/internals/domain/service"
	logger "project/package/utils/pkg"

	"go.uber.org/zap"
)

type ParseFileUseCase struct {
	// SqlRepo        repository.ISqlRepository
	Service service.Service
	Logger  *zap.SugaredLogger
}

func NewParseFileUseCase(service service.Service) *ParseFileUseCase {
	return &ParseFileUseCase{
		Service: service,
		Logger:  logger.Logger,
	}
}

func (uc *ParseFileUseCase) GenerateCertificateHTML(templatePath string, certificateData entity.CertificateData) (string, error) {

	htmlContent, err := uc.Service.ParseAndExecute(templatePath, certificateData)
	if err != nil {
		uc.Logger.Errorw("[certificate_usecase] Failed to generate HTML", "error", err)
		return "", err
	}

	uc.Logger.Infow("[certificate_usecase] Successfully generated certificate HTML")
	return htmlContent, nil
}
