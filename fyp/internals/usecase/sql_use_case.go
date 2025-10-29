package usecase

import (
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	"project/internals/domain/service"
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
		PublicKey:             "", //! public key is empty
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

	institutionInfo, er := uc.SqlRepo.GetPendingInstitutionFromInstitutionID(faculty.InstitutionID)
	if er != nil {
		return "", nil, er
	}
	return facultyID, institutionInfo, nil
}

// func (uc *SqlUseCase) GetInstitutionsToBeVerifiedUseCase() ([]entity.Institution, error) {
// 	return uc.SqlRepo.GetToBeVerifiedInstitutions()
// }
