package filehandling

import (
	"project/internals/domain/entity"
	"project/internals/usecase"
	err "project/package/errors"
	"project/package/utils/common"
)

type Controller struct {
	ParseFileUseCase  *usecase.ParseFileUseCase
	BlockChainUseCase *usecase.BlockChainUseCase
}

func NewController(parseFileUseCase *usecase.ParseFileUseCase, blockChainUseCase *usecase.BlockChainUseCase) *Controller {
	return &Controller{ParseFileUseCase: parseFileUseCase, BlockChainUseCase: blockChainUseCase}
}

func (c *Controller) HandleGetHTMLFile(request map[string]string) entity.Response {
	//! extract type of certificate in future implementation

	templatePath := "../internals/templates/certificate.html"
	fakeCertificateData, er := c.BlockChainUseCase.GetCertificateData()
	if er != nil {
		return common.HandleErrorResponse(500, err.ErrCreatingInstitutionFacultyString, er)
	}
	htmlString, er := c.ParseFileUseCase.GenerateCertificateHTML(templatePath, fakeCertificateData)
	if er != nil {
		return common.HandleErrorResponse(500, err.ErrParsingFileString, er)
	}

	return common.HandleSuccessResponse(htmlString)

}
