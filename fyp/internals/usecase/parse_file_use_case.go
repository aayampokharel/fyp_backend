package usecase

import (
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	"project/internals/domain/service"
)

type ParseFileUseCase struct {
	SqlRepo repository.ISqlRepository
	Service service.Service
}

func NewParseFileUseCase(service service.Service) *ParseFileUseCase {
	return &ParseFileUseCase{
		Service: service,
	}
}

func (uc *ParseFileUseCase) GenerateCertificateHTML(id, url, templatePath string, certificateData entity.CertificateData) (string, error) {
	qrCodeBase64, er := uc.Service.GenerateQRCodeBase64(id, url)
	if er != nil {
		return "", er
	}

	certificateDataWithQR := entity.CertificateDataWithQRCode{
		CertificateData: certificateData,
		QRCodeBase64:    qrCodeBase64,
	}

	htmlContent, err := uc.Service.ParseAndExecute(templatePath, certificateDataWithQR)
	if err != nil {
		uc.Service.Logger.Errorw("[certificate_usecase] Failed to generate HTML", "error", err)
		return "", err
	}

	uc.Service.Logger.Infow("[certificate_usecase] Successfully generated certificate HTML")
	return htmlContent, nil
}

func (uc *ParseFileUseCase) GenerateAndStoreCertificatePDF(htmlContent string, pdfFileEntity entity.PDFFileEntity) error {
	pdfBytes, er := uc.Service.ConvertHTMLToPDF(htmlContent)
	if er != nil {
		uc.Service.Logger.Errorln("[certificate_usecase] error while generating pdfbytes ", er)
		return er
	}
	pdfFileEntity.PDFData = pdfBytes

	er = uc.SqlRepo.InsertPDFFile(pdfFileEntity)
	if er != nil {
		uc.Service.Logger.Errorln("[certificate_usecase] error while storing pdfbytes ", er)
		return er
	}

	return nil
}

func (uc *ParseFileUseCase) RetrievePDFFileByFileIDOrCategoryID(pdfFileId string, categoryID string, isDownloadAll bool) ([]entity.PDFFileEntity, error) {
	return uc.SqlRepo.RetrievePDFFileByFileIDOrCategoryID(pdfFileId, categoryID, isDownloadAll)

}
