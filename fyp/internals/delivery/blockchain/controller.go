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
	var basicStudentInfoDto []BasicStudentInfoDto
	blockChainLength := c.useCase.GetBlockChainLength()
	if blockChainLength == 0 {
		//mock data
		if er := c.useCase.InsertGenesisBlock(); er != nil {
			log.Println(er)
			return common.HandleErrorResponse(401, er.Error(), er)
		}
	}
	// certificateData, er := c.useCase.GetCertificateData()
	// if er != nil {
	// 	log.Println(er)
	// 	return nil, er
	// }
	//
	//
	//
	// pdfFileCategory, er := request.ToPdfFileCategoryEntity()
	// if er != nil {
	// 	log.Println(er)
	// 	return common.HandleErrorResponse(401, er.Error(), er)
	// }
	// insertedpdfFileCategory, er := c.useCase.SqlRepo.InsertAndGetPDFCategory(pdfFileCategory)
	// if er != nil {
	// 	return common.HandleErrorResponse(401, er.Error(), er)
	// }
	// c.useCase.SqlRepo.InsertAndGetPDFCategory(pdfFileCategory)

	for i := 0; i < len(request.CertificateData); i++ {

		certificateData := request.CertificateData[i].ToEntity()

		latestBlockFromChain, newBlock, latestBlockFromChainCertificateLength, newBlockCertificateLength, er := c.useCase.CompleteBlockFromCertificate(certificateData)
		if er != nil {
			log.Println(er)
			return common.HandleErrorResponse(401, er.Error(), er)
		}
		var strNodeInfoMap map[int]string
		strNodeInfoMap, er = c.useCase.BroadcastNewBlock(newBlock)
		if er != nil {
			log.Println(er)
			return common.HandleErrorResponse(401, er.Error(), er)
		}
		log.Println("Acknowledgement from nodes: ", strNodeInfoMap)

		if er = c.useCase.UpsertBlockChain(*latestBlockFromChain, *newBlock, latestBlockFromChainCertificateLength, newBlockCertificateLength); er != nil {
			log.Println(er)
			return common.HandleErrorResponse(401, er.Error(), er)
		}

		templatePath := constants.TemplateBasePath + constants.CertificateTemplate
		//fakeCertificateData, er := c.BlockChainUseCase.GetCertificateData()
		// if er != nil {
		// 	return common.HandleFileErrorResponse(500, er.ErrCreatingInstitutionFacultyString, er)
		// }
		htmlString, er := c.ParseFileUseCase.GenerateCertificateHTML("123", "url", templatePath, *certificateData)
		if er != nil {
			return common.HandleErrorResponse(500, err.ErrParsingFileString, er)
		}
		//c.sqlUseCase.Service.Logger.Infoln("[certificate_usecase] htmlString", htmlString)
		pdfBytes, er := c.ParseFileUseCase.GenerateAndGetCertificatePDF(htmlString)
		if er != nil {
			return common.HandleErrorResponse(500, err.ErrParsingFileString, er)
		}

		pdfEntityWithoutData := FromPDFFileCategoryToPDFFileEntity(request.CategoryID, certificateData.StudentName, request.InstitutionFacultyName, i)

		pdfEntityWithoutData.PDFData = pdfBytes
		fileID, er := c.sqlUseCase.InsertPDFFileUseCase(pdfEntityWithoutData)
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

		basicStudentInfoDto = append(basicStudentInfoDto, BasicStudentInfoDto{
			StudentID:   certificateData.StudentID,
			StudentName: certificateData.StudentName,
			FileID:      fileID,
			FileName:    pdfEntityWithoutData.FileName,
			FacultyName: request.InstitutionFacultyName,
		})

	}
	return common.HandleSuccessResponse(CreateAllCertificateResponse{
		Message:     "All Certificates Inserted Successfully",
		StudentList: basicStudentInfoDto,
	})
}
