package filehandling

import (
	"project/constants"
	"project/internals/domain/entity"
	"project/internals/usecase"
	"project/package/enum"
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

// func (c *Controller) HandleCreatePDFFile(request )
//creation of pdf should be handled from blockchain package .

func (c *Controller) HandleGetHTMLFile(request GetRequestQueryType) entity.Response {

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

func (c *Controller) HandleGetPDFFileInList(request GetRequestQueryType) entity.Response {
	var responseWithFileTypeAndCount ResponseWithFileTypeAndCount
	checkedMap := common.CheckMapKeysReturnValues(request, "category_id", "file_id", "is_download_all")
	if checkedMap == nil {
		return common.HandleErrorResponse(500, err.ErrParsingQueryParametersString, nil)
	}
	categoryID := checkedMap["category_id"]
	fileID := checkedMap["file_id"]
	isDownloadAll, er := common.CheckBoolFromString(checkedMap["is_download_all"])
	if er != nil {
		return common.HandleErrorResponse(500, err.ErrDataTypeMismatchString, er)
	}
	pdfFileEntity, er := c.ParseFileUseCase.RetrievePDFFileByFileIDOrCategoryID(fileID, categoryID, isDownloadAll)
	if er != nil {
		return common.HandleErrorResponse(500, err.ErrParsingFileString, er)
	}

	responseWithFileTypeAndCount.Count = len(pdfFileEntity)
	responseWithFileTypeAndCount.Data = pdfFileEntity
	if isDownloadAll && len(pdfFileEntity) > 1 {
		responseWithFileTypeAndCount.FileType = enum.ZIP

		return common.HandleSuccessResponse(responseWithFileTypeAndCount)
	}

	responseWithFileTypeAndCount.FileType = enum.PDF
	return common.HandleSuccessResponse(responseWithFileTypeAndCount)

}
