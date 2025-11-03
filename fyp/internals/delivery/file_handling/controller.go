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

func (c *Controller) HandleGetHTMLFile(request map[string]string) entity.FileResponse {

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
		return common.HandleFileErrorResponse(500, err.ErrCreatingInstitutionFacultyString, er)
	}
	htmlString, er := c.ParseFileUseCase.GenerateCertificateHTML("123", "url", templatePath, *fakeCertificateData)
	if er != nil {
		return common.HandleFileErrorResponse(500, err.ErrParsingFileString, er)
	}

	return common.HandleFileSuccessResponse(enum.HTML, "", []byte(htmlString))

}

func (c *Controller) HandleGetPDFFileInList(request map[string]string) entity.FileResponse {
	var fileName string
	checkedMap, er := common.CheckMapKeysReturnValues(request, GetPDFFileInListQuery)
	if er != nil {
		return common.HandleFileErrorResponse(500, err.ErrParsingQueryParametersString, nil)
	}

	categoryID := checkedMap[CategoryId]
	fileID := checkedMap[FileID]
	categoryName := checkedMap[CategoryName]
	isDownloadAll, er := common.ConvertToBool(checkedMap[IsDownloadAll])
	if er != nil {
		return common.HandleFileErrorResponse(500, err.ErrDataTypeMismatchString, er)
	}
	pdfFileEntity, er := c.ParseFileUseCase.RetrievePDFFileByFileIDOrCategoryID(fileID, categoryID, isDownloadAll)
	if er != nil {
		return common.HandleFileErrorResponse(500, err.ErrParsingFileString, er)
	}
	if pdfFileEntity == nil {
		return common.HandleFileErrorResponse(400, err.ErrFileNotFoundString, nil)
	}

	if isDownloadAll && len(pdfFileEntity) > 1 {
		fileName = categoryName + "_" + common.GenerateUUID(6)
		zipBytes, er := c.ParseFileUseCase.Service.CreateZipUsingPDF(pdfFileEntity)
		if er != nil {
			return common.HandleFileErrorResponse(500, err.ErrZipWritingString, er)
		}
		return common.HandleFileSuccessResponse(enum.ZIP, fileName, zipBytes)
	}

	fileName = pdfFileEntity[0].FileName
	if pdfFileEntity[0].PDFData == nil {
		return common.HandleFileErrorResponse(500, err.ErrParsingFileString, er)
	}
	c.ParseFileUseCase.Service.Logger.Debugln(pdfFileEntity[0].PDFData[:1000])

	fileName = pdfFileEntity[0].CategoryID + "_" + fileName

	return common.HandleFileSuccessResponse(enum.PDF, fileName, pdfFileEntity[0].PDFData)

}
