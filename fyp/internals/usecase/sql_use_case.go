package usecase

import (
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	"project/internals/domain/service"
	logger "project/package/utils/pkg"

	"go.uber.org/zap"
)

type SqlUseCase struct {
	SqlRepo repository.ISqlRepository
	Service service.Service
	Logger  *zap.SugaredLogger
}

func NewSqlUseCase(sqlRepo repository.ISqlRepository, service service.Service) *SqlUseCase {
	return &SqlUseCase{
		SqlRepo: sqlRepo,
		Service: service,
		Logger:  logger.Logger,
	}
}

func (uc *SqlUseCase) InsertUserAccountUseCase(userAccount entity.UserAccount, institutionId string, institutionLogoBase64 string) error {
	er := uc.SqlRepo.InsertUserAccounts(userAccount)
	if er != nil {
		return er
	}
	insUser := entity.InstitutionUser{
		UserID:                userAccount.ID,
		InstitutionID:         institutionId,
		InstitutionLogoBase64: institutionLogoBase64,
		PublicKey:             "", //! public key is empty
	}
	return uc.SqlRepo.InsertInstitutionUser(insUser)

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

func (uc *SqlUseCase) InsertFacultyUseCase(faculty entity.InstitutionFaculty, institutionID string) error {
	return uc.SqlRepo.InsertFaculty(faculty, institutionID)
}
