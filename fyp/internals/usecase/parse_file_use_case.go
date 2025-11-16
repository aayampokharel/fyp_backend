package usecase

import (
	"project/constants"
	"project/internals/data/config"
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	"project/internals/domain/service"
)

type ParseFileUseCase struct {
	env     *config.Env
	SqlRepo repository.ISqlRepository
	Service service.Service
}

func NewParseFileUseCase(service service.Service, env *config.Env, sqlRepo repository.ISqlRepository) *ParseFileUseCase {
	return &ParseFileUseCase{
		Service: service,
		SqlRepo: sqlRepo,
		env:     env,
	}
}

func (uc *ParseFileUseCase) GenerateCertificateHTML(id, hash, url, templatePath string, certificateData entity.CertificateData, institutionLogo, authorityNameWithSignature string) (string, error) {
	qrUrl := uc.env.GetValueForKey(constants.PinggyQrUrl)
	qrCodeBase64, er := uc.Service.GenerateQRCodeBase64(id, hash, qrUrl)
	if er != nil {
		return "", er
	}

	authorityEntityList, er := entity.AuthorityWithSignatureEntity{}.FromString(authorityNameWithSignature)
	if er != nil {
		return "", er
	}
	certificateDataWithLogosAndQR := entity.CertificateDataWithLogosAndQRCode{
		CertificateDataWithLogos: entity.CertificateDataWithLogos{
			CertificateData:                certificateData,
			InstitutionLogoBase64:          institutionLogo,
			AuthorityWithSignatureEntities: authorityEntityList,
		},
		QRCodeBase64: qrCodeBase64,
	}

	htmlContent, err := uc.Service.ParseAndExecute(templatePath, certificateDataWithLogosAndQR)
	if err != nil {
		uc.Service.Logger.Errorw("[certificate_usecase] Failed to generate HTML", "error", err)
		return "", err
	}

	uc.Service.Logger.Infow("[certificate_usecase] Successfully generated certificate HTML")
	return htmlContent, nil
}

func (uc *ParseFileUseCase) GenerateAndGetCertificatePDF(htmlContent string) ([]byte, error) {
	// uc.Service.Logger.Debugln(htmlContent)
	pdfBytes, er := uc.Service.ConvertHTMLToPDF(htmlContent)
	if er != nil {
		uc.Service.Logger.Errorln("[certificate_usecase] error while generating pdfbytes ", er)
		return nil, er
	}

	return pdfBytes, nil
}

func (uc *ParseFileUseCase) RetrievePDFFileByFileIDOrCategoryID(pdfFileId string, categoryID string, isDownloadAll bool) ([]entity.PDFFileEntity, error) {
	return uc.SqlRepo.RetrievePDFFileByFileIDOrCategoryID(pdfFileId, categoryID, isDownloadAll)

}
