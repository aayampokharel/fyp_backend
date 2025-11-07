package filehandling

import (
	"encoding/base64"
	"project/constants"
	"project/internals/domain/entity"
	"project/internals/usecase"
	"project/package/enum"
	err "project/package/errors"
	"project/package/utils/common"
	"strings"
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
	c.BlockChainUseCase.Service.Logger.Infoln("[handle_get_pdf_file_in_list] Info: HandleGetPDFFileInList::", request)
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
	//! I have to include principal signature in certificate as well .
	if isDownloadAll && len(pdfFileEntity) > 1 {
		fileName = common.GenerateFileNameWithExtension(categoryName, 6, "zip")
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

func (c *Controller) HandleGetImageFile(request GetImageFileRequestDto) entity.FileResponse {
	cleanBase64 := request.ImageBase64
	var cleanBase64slice []string
	var encodedString string
	if strings.Contains(cleanBase64, "base64,") {
		cleanBase64slice = strings.Split(cleanBase64, ",")
		if len(cleanBase64slice) != 2 {
			return common.HandleFileErrorResponse(400, err.ErrInvalidBase64.Error(), nil)
		}
		encodedString = cleanBase64slice[1]
	} else {
		encodedString = cleanBase64
	}
	decodedBytes, er := base64.StdEncoding.DecodeString(encodedString)
	if er != nil {
		return common.HandleFileErrorResponse(500, err.ErrParsingFileString, er)
	}
	if len(decodedBytes) == 0 {
		return common.HandleFileErrorResponse(400, err.ErrInvalidLengthString, nil)
	}
	resultImageBytes, er := c.ParseFileUseCase.Service.RemoveBackgroundService(decodedBytes)
	if er != nil {
		return common.HandleFileErrorResponse(500, err.ErrParsingFileString, er)
	}
	return common.HandleFileSuccessResponse(enum.IMAGE, "removed_background_"+request.ImageName, resultImageBytes)

}
