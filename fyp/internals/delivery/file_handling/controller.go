package filehandling

import (
	"encoding/base64"
	"encoding/hex"
	"project/constants"
	"project/internals/domain/entity"
	"project/internals/usecase"
	"project/package/enum"
	err "project/package/errors"
	"project/package/utils/common"
	"strings"
	"time"
)

type Controller struct {
	PinggyQrUrl          string
	currentMappedTCPPort int
	ParseFileUseCase     *usecase.ParseFileUseCase
	SqlUseCase           *usecase.SqlUseCase
	BlockChainUseCase    *usecase.BlockChainUseCase
	PbftUseCase          *usecase.PBFTUseCase
}

func NewController(parseFileUseCase *usecase.ParseFileUseCase, blockChainUseCase *usecase.BlockChainUseCase, currentMappedTcpPort int, pbftUseCase *usecase.PBFTUseCase, pingyUrl string, sqlUseCase *usecase.SqlUseCase) *Controller {
	return &Controller{ParseFileUseCase: parseFileUseCase, BlockChainUseCase: blockChainUseCase, currentMappedTCPPort: currentMappedTcpPort, PbftUseCase: pbftUseCase, PinggyQrUrl: pingyUrl, SqlUseCase: sqlUseCase}
}

func (c *Controller) HandleGetHTMLFile(request map[string]string) entity.FileResponse {
	c.BlockChainUseCase.Service.Logger.Debugln("currentPORT selected:", c.currentMappedTCPPort)

	requestMap, er := common.CheckMapKeysReturnValues(request, GetHTMLRequestQuery)
	if er != nil {
		return common.HandleFileErrorResponse(400, err.ErrParsingQueryParametersString, er)
	}

	hashString := requestMap[CertificateHash]
	hashBytes, er := hex.DecodeString(hashString)
	if er != nil {
		return common.HandleFileErrorResponse(500, err.ErrParsingCertificateIDString, er)
	}
	certificateID := requestMap[CertificateID]
	pbftExecuteResult := make(chan entity.PBFTExecutionResultEntity)
	c.PbftUseCase.SendPBFTMessageToPeer(entity.PBFTMessage{
		VerificationType: enum.INITIAL,
		QRVerificationRequestData: entity.QRVerificationRequestData{
			CertificateHash: hashBytes,
			CertificateID:   certificateID,
		},
	}, c.currentMappedTCPPort, pbftExecuteResult)

	select {
	case result := <-pbftExecuteResult:
		if result.Result {

			resultCertificateData, er := c.BlockChainUseCase.GetCertificateDataByCertificateIDAndHashUseCase(hashString, certificateID)
			if er != nil {
				return common.HandleFileErrorResponse(500, err.ErrParsingFileString, er)
			}
			institutionID := resultCertificateData.InstitutionID
			institutionFacultyID := resultCertificateData.InstitutionFacultyID
			// institutionID, institutionFacultyID, er := c.BlockChainUseCase.GetInstitutionIDAndFacultyIDFromCertificateIDAndHashUseCase(certificateID, hashString)

			institutionLogo, authorityNameWithSignature, er := c.SqlUseCase.GetAllLogosForCertificateUseCase(institutionID, institutionFacultyID)
			if er != nil {
				return common.HandleFileErrorResponse(500, err.ErrParsingFileString, er)
			}

			// fakeCertificateData, _ := c.BlockChainUseCase.GetCertificateData()
			templatePath := constants.TemplateBasePath + constants.CertificateTemplate
			qrUrl := c.PinggyQrUrl
			htmlString, er := c.ParseFileUseCase.GenerateCertificateHTML(certificateID, hashString, qrUrl, templatePath, resultCertificateData, institutionLogo, authorityNameWithSignature)
			if er != nil {
				return common.HandleFileErrorResponse(500, err.ErrParsingFileString, er)
			}

			return common.HandleFileSuccessResponse(enum.HTML, "certificate.html", []byte(htmlString))

		}

	case <-time.After(20 * time.Second):
		return common.HandleFileErrorResponse(500, err.ErrParsingFileString, nil)
	}
	return common.HandleFileErrorResponse(500, err.ErrParsingFileString, nil)

	// if er != nil {
	// 	return common.HandleFileErrorResponse(500, err.ErrParsingFileString, er)
	// }

	// return common.HandleFileSuccessResponse(enum.HTML, "", []byte("helllo change this string to variable "))

}

func (c *Controller) HandleGetPDFFileInList(request map[string]string) entity.FileResponse {
	var fileName string
	c.BlockChainUseCase.Service.Logger.Infoln("[handle_get_pdf_file_in_list] Info: HandleGetPDFFileInList::", request)
	checkedMap, er := common.CheckMapKeysReturnValues(request, GetPDFFileInListQuery)
	if er != nil {
		return common.HandleFileErrorResponse(400, err.ErrParsingQueryParametersString, nil)
	}

	categoryID := checkedMap[CategoryId]
	fileID := checkedMap[FileID]
	categoryName := checkedMap[CategoryName]
	isDownloadAll, er := common.ConvertToBool(checkedMap[IsDownloadAll])
	if er != nil {
		return common.HandleFileErrorResponse(400, err.ErrDataTypeMismatchString, er)
	}
	pdfFileEntity, er := c.ParseFileUseCase.RetrievePDFFileByFileIDOrCategoryID(fileID, categoryID, isDownloadAll)
	if er != nil {
		return common.HandleFileErrorResponse(500, err.ErrParsingFileString, er)
	}
	if pdfFileEntity == nil {
		return common.HandleFileErrorResponse(404, err.ErrFileNotFoundString, nil)
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
		return common.HandleFileErrorResponse(422, err.ErrParsingFileString, er)
	}
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
		return common.HandleFileErrorResponse(422, err.ErrParsingFileString, er)
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
