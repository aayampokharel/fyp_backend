package repository

import (
	"project/internals/domain/entity"
	"project/package/enum"
	"time"
)

type ISqlRepository interface {
	CheckDuplicationByInstitutionInfo(institution entity.Institution) error
	GetUserAccountsForEmail(userEmail string) (entity.UserAccount, error)
	InsertInstitutions(institution entity.Institution) (string, error)
	InsertUserAccounts(userAccounts entity.UserAccount) (string, error)
	InsertInstitutionUser(institutionUser entity.InstitutionUser) error
	InsertFaculty(faculty entity.InstitutionFaculty) (string, error)
	GetUserIDByInstitutionID(institutionID string) (string, error)
	UpdateSignUpCompletedByInstitutionID(institutionID string) error
	UpdateIsActiveByInstitutionID(institutionID string, isActive bool) error
	GetAllPendingInstitutionsForAdmin(adminID string) ([]entity.Institution, error)
	//GetPendingInstitutionFromInstitutionID(institutionID string) (*entity.Institution, error)
	GetPDFCategoriesList(institutionID, institutionFacultyID string) ([]entity.PDFFileCategoryEntity, error)

	VerifyRoleLogin(userMail, password string, role enum.ROLE) (string, time.Time, error)
	InsertPDFFile(pdfFile entity.PDFFileEntity) error
	InsertAndGetPDFCategory(pdfFileCategory entity.PDFFileCategoryEntity) (*entity.PDFFileCategoryEntity, error)
	RetrievePDFFileByFileIDOrCategoryID(fileID string, categoryID string, isDownloadAll bool) ([]entity.PDFFileEntity, error)
	InsertBlockWithSingleCertificate(blockHeader entity.Header, certificateData entity.CertificateData, certificatePositionZeroIndex int) error
	UpdateBlockDataByID(blockHeader entity.Header, id string) error
	InsertCertificate(certificate entity.CertificateData, blockNumber int, certificatePositionZeroIndex int) error
	GetFacultyPublicKey(id string) (string, error)
	GetInfoFromPdfFilesCategories(categoryID string) (*entity.PDFFileCategoryEntity, error)
	GetInstitutionInfoFromInstitutionID(institutionID string) (*entity.Institution, error)
	GetInstitutionsForUser(userID string) ([]entity.Institution, error)
}
