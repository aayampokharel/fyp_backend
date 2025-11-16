package usecase

import (
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	"project/internals/domain/service"
	"project/package/enum"
	err "project/package/errors"
	"project/package/utils/common"
	"time"
)

type SqlUseCase struct {
	SqlRepo repository.ISqlRepository
	Service service.Service
}

func NewSqlUseCase(sqlRepo repository.ISqlRepository, service service.Service) *SqlUseCase {
	return &SqlUseCase{
		SqlRepo: sqlRepo,
		Service: service,
	}
}

func (uc *SqlUseCase) InsertUserAccountUseCase(userAccount entity.UserAccount, institutionId string, institutionLogoBase64 string) (createdAt string, userAccountID string, er error) {
	createdAt, er = uc.SqlRepo.InsertUserAccounts(userAccount)
	if er != nil {
		return "", "", er
	}
	insUser := entity.InstitutionUser{
		UserID:                userAccount.ID,
		InstitutionID:         institutionId,
		InstitutionLogoBase64: institutionLogoBase64,
	}
	return createdAt, userAccount.ID, uc.SqlRepo.InsertInstitutionUser(insUser)

}

func (uc *SqlUseCase) CheckDuplicationByInstitutionInfoUseCase(institution entity.Institution) error {
	return uc.SqlRepo.CheckDuplicationByInstitutionInfo(institution)
}

func (uc *SqlUseCase) InsertInstitutionsUseCase(institution entity.Institution) (string, error) {
	er := uc.SqlRepo.CheckDuplicationByInstitutionInfo(institution)
	if er != nil {
		return "", er
	}

	return uc.SqlRepo.InsertInstitutions(institution)
}

func (uc *SqlUseCase) InsertFacultyAndRetrieveInstitutionUseCase(faculty entity.InstitutionFaculty) (string, *entity.Institution, error) {
	facultyID, er := uc.SqlRepo.InsertFaculty(faculty)
	if er != nil {
		return "", nil, er
	}
	er = uc.SqlRepo.UpdateSignUpCompletedByInstitutionID(faculty.InstitutionID)
	if er != nil {
		return "", nil, er
	}

	institutionInfo, er := uc.SqlRepo.GetInstitutionInfoFromInstitutionID(faculty.InstitutionID)
	if er != nil {
		return "", nil, er
	}
	if institutionInfo.IsActive != nil {
		return facultyID, nil, err.ErrInstitutionAlreadyVerified
	}
	return facultyID, institutionInfo, nil
}

func (uc *SqlUseCase) InsertAndGetPDFCategoryUseCase(pdfCategory entity.PDFFileCategoryEntity) (*entity.PDFFileCategoryEntity, error) {
	return uc.SqlRepo.InsertAndGetPDFCategory(pdfCategory)
}

func (uc *SqlUseCase) InsertPDFFileUseCase(pdfFile entity.PDFFileEntity) (string, error) {
	er := uc.SqlRepo.InsertPDFFile(pdfFile)

	if er != nil {
		return "", er
	}
	return pdfFile.FileID, nil

}

// func (uc *SqlUseCase) GetInstitutionsToBeVerifiedUseCase() ([]entity.Institution, error) {
// 	return uc.SqlRepo.GetToBeVerifiedInstitutions()
// }

func (uc *SqlUseCase) GetInstitutionInfoUseCase(institutionID string) (*entity.Institution, error) {
	return uc.SqlRepo.GetInstitutionInfoFromInstitutionID(institutionID)
}

func (uc *SqlUseCase) GetAllPendingInstitutionsUseCase(adminID string) ([]entity.Institution, error) {
	return uc.SqlRepo.GetAllPendingInstitutionsForAdmin(adminID)
}
func (uc *SqlUseCase) GetPDFCategoriesListUseCase(institutionID, institutionFacultyID string) ([]entity.PDFFileCategoryEntity, error) {
	return uc.SqlRepo.GetPDFCategoriesList(institutionID, institutionFacultyID)
}

func (uc *SqlUseCase) GetAdminDetailsUseCase(userID string, role enum.ROLE) (*entity.AdminDashboardCountsEntity, []entity.Institution, error) {

	adminDashboardDetails, er := uc.SqlRepo.GetAdminDashboardCounts(userID)
	if er != nil {
		return nil, nil, er
	}
	pendingInstitutions, er := uc.SqlRepo.GetAllPendingInstitutionsForAdmin(userID)
	if er != nil {
		return nil, nil, er
	}

	return adminDashboardDetails, pendingInstitutions, nil
}

func (uc *SqlUseCase) VerifyUserLoginUseCase(userEmail, password string, role enum.ROLE) (string, time.Time, error) {
	hashedPassword, _, er := common.HashData(password)
	if er != nil {
		return "", time.Time{}, er
	}
	userID, createdAt, er := uc.SqlRepo.VerifyRoleLogin(userEmail, hashedPassword, role)
	if er != nil {
		return "", time.Time{}, er
	}

	return userID, createdAt, nil
}
func (uc *SqlUseCase) GetInstitutionsForUserUseCase(userID string) ([]entity.Institution, error) {
	return uc.SqlRepo.GetInstitutionsForUser(userID)
}
func (uc *SqlUseCase) GetFacultyListUseCase(institutionID string) ([]entity.InstitutionFaculty, error) {
	return uc.SqlRepo.GetFacultyListInfoFromInstitutionID(institutionID)
}
func (uc *SqlUseCase) GetAllLogosForCertificateUseCase(institutionID, facultyID string) (string, string, error) {
	return uc.SqlRepo.GetAllLogosForCertificate(institutionID, facultyID)
}
func (uc *SqlUseCase) GetFacultiesForInstitutionIDUseCase(institutionID string) ([]entity.InstitutionFaculty, error) {
	return uc.SqlRepo.GetFacultiesForInstitutionID(institutionID)
}

func (uc *SqlUseCase) UpdateIsActiveByInstitutionIDUseCase(institutionID string, isActive bool) error {
	return uc.SqlRepo.UpdateIsActiveByInstitutionID(institutionID, isActive)
}
func (uc *SqlUseCase) DeleteUserByUserIDUseCase(userID string) (string, error) {
	return uc.SqlRepo.DeleteUserByUserID(userID)
}

func (uc *SqlUseCase) InsertBlockWithFullCertificates(blockHeader entity.Header, certificates [4]entity.CertificateData) error {
	return uc.SqlRepo.InsertBlockWithAllCertificates(blockHeader, certificates)
}
