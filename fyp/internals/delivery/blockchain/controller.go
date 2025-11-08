package delivery

import (
	"log"
	"project/constants"
	"project/internals/domain/entity"
	"project/internals/usecase"
	err "project/package/errors"
	"project/package/utils/common"
)

type Controller struct {
	useCase          usecase.BlockChainUseCase
	ParseFileUseCase *usecase.ParseFileUseCase
	sqlUseCase       *usecase.SqlUseCase
}

func NewController(useCase usecase.BlockChainUseCase, parseFileUseCase *usecase.ParseFileUseCase, sqlUseCase *usecase.SqlUseCase) *Controller {
	return &Controller{useCase: useCase, ParseFileUseCase: parseFileUseCase, sqlUseCase: sqlUseCase}
}

func (c *Controller) InsertNewCertificateData(request CreateCertificateDataRequest) entity.Response {
	//var basicStudentInfoDto []BasicStudentInfoDto
	blockChainLength := c.useCase.GetBlockChainLength()
	if blockChainLength == 0 {
		//mock data
		if er := c.useCase.InsertGenesisBlock(); er != nil {
			log.Println(er)
			return common.HandleErrorResponse(500, er.Error(), er)
		}
	}

	for i := 0; i < len(request.CertificateData); i++ {

		certificateData, er := request.CertificateData[i].ToEntity(request.CategoryID)
		if er != nil {
			log.Println(er)
			return common.HandleErrorResponse(400, er.Error(), er)
		}

		latestBlockFromChain, newBlock, latestBlockFromChainCertificateLength, newBlockCertificateLength, er := c.useCase.CompleteBlockFromCertificate(certificateData)
		if er != nil {
			log.Println(er)
			return common.HandleErrorResponse(500, er.Error(), er)
		}
		var strNodeInfoMap map[int]string
		strNodeInfoMap, er = c.useCase.BroadcastNewBlock(newBlock)
		if er != nil {
			log.Println(er)
			return common.HandleErrorResponse(500, er.Error(), er)
		}
		log.Println("Acknowledgement from nodes: ", strNodeInfoMap)

		if er = c.useCase.UpsertBlockChain(*latestBlockFromChain, *newBlock, latestBlockFromChainCertificateLength, newBlockCertificateLength); er != nil {
			log.Println(er)
			return common.HandleErrorResponse(500, er.Error(), er)
		}

		templatePath := constants.TemplateBasePath + constants.CertificateTemplate
		//fakeCertificateData, er := c.BlockChainUseCase.GetCertificateData()
		// if er != nil {
		// 	return common.HandleFileErrorResponse(500, er.ErrCreatingInstitutionFacultyString, er)
		// }
		htmlString, er := c.ParseFileUseCase.GenerateCertificateHTML("123", "url", templatePath, *certificateData)
		if er != nil {
			return common.HandleErrorResponse(422, err.ErrParsingFileString, er)
		}
		//c.sqlUseCase.Service.Logger.Infoln("[certificate_usecase] htmlString", htmlString)
		pdfBytes, er := c.ParseFileUseCase.GenerateAndGetCertificatePDF(htmlString)
		if er != nil {
			return common.HandleErrorResponse(422, err.ErrParsingFileString, er)
		}

		pdfEntityWithoutData := FromPDFFileCategoryToPDFFileEntity(request.CategoryID, certificateData.StudentName, request.InstitutionFacultyName, certificateData.PDFFileID, i)

		pdfEntityWithoutData.PDFData = pdfBytes
		_, er = c.sqlUseCase.InsertPDFFileUseCase(pdfEntityWithoutData)
		if er != nil {
			c.useCase.Service.Logger.Errorln("[certificate_usecase] error while storing pdfbytes ", er)
			return common.HandleErrorResponse(500, err.ErrCreatingInstitutionFacultyString, er)
		}

		// c.useCase.SqlRepo.InsertPDFFile(entity.PDFFileEntity{
		// 	FileID:     common.GenerateUUID(16),
		// 	CategoryID: insertedpdfFileCategory.CategoryID,
		// 	FileName:   certificateData.FileName,
		// 	PDFData:    certificateData.PDFData,
		// })
		// blockchain, _ := c.useCase.BlockChainRepo.GetBlockChain()

		// basicStudentInfoDto = append(basicStudentInfoDto, BasicStudentInfoDto{
		// 	StudentID:   certificateData.StudentID,
		// 	StudentName: certificateData.StudentName,
		// 	FileID:      fileID,
		// 	FileName:    pdfEntityWithoutData.FileName,
		// 	FacultyName: request.InstitutionFacultyName,
		// })

	}
	return common.HandleSuccessResponse(CreateAllCertificateResponse{
		Message: "All Certificates Inserted Successfully",
	})
}

func (c *Controller) GetCertificateDataList(request map[string]string) entity.Response {
	requestMap, er := common.CheckMapKeysReturnValues(request, GetCertificateDataListRequestQuery)
	if er != nil {
		return common.HandleErrorResponse(400, err.ErrParsingQueryParametersString, er)
	}
	institutionID := requestMap[InstitutionID]
	institutionFacultyID := requestMap[InstitutionFacultyID]
	categoryID := requestMap[CategoryID]
	certificates, er := c.useCase.GetCertificateDataListUseCase(institutionID, institutionFacultyID, categoryID)
	if er != nil {
		return common.HandleErrorResponse(500, er.Error(), er)
	}
	return common.HandleSuccessResponse(certificates)
}
