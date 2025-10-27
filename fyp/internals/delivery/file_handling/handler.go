package filehandling

import (
	"project/constants"
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

	//// in preview dont show this :
	//  <div class="block-info">
	//     Block: {{.BlockNumber}}<br>
	//     Position: {{.Position}}
	// </div>
	//===============================================
	// 	<div class="verification-info">
	//   <h3>Blockchain Verification Details</h3>
	//   <p><strong>Block Number:</strong> {{.BlockNumber}}</p>
	//   <p><strong>Transaction Index:</strong> {{.Position}}</p>
	//   <p><strong>Data Hash:</strong> {{.DataHash}}</p>
	// </div>
	//// but for /verify/ show this as well not for preview/download feature.

	//! extract type of certificate in future implementation
	//// add college seal section in database .
	// templatePath := "../internals/templates/certificate.html"
	templatePath := constants.TemplateBasePath + constants.CertificateTemplate
	fakeCertificateData, er := c.BlockChainUseCase.GetCertificateData()
	if er != nil {
		return common.HandleErrorResponse(500, err.ErrCreatingInstitutionFacultyString, er)
	}
	htmlString, er := c.ParseFileUseCase.GenerateCertificateHTML("123", "url", templatePath, fakeCertificateData)
	if er != nil {
		return common.HandleErrorResponse(500, err.ErrParsingFileString, er)
	}

	return common.HandleSuccessResponse(htmlString)

}
