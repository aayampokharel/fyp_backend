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

func (uc *SqlUseCase) InsertUserAccountUseCase(userAccount entity.UserAccount) error {
	return uc.SqlRepo.InsertUserAccounts(userAccount)
}

func (uc *SqlUseCase) InsertInstitutionsUseCase(institution entity.Institution) error {
	return uc.SqlRepo.InsertInstitutions(institution)
}

func (uc *SqlUseCase) InsertInstitutionUsersUseCase(institutionUser entity.InstitutionUser, institutionID, userID string) error {
	return uc.SqlRepo.InsertInstitutionUser(institutionUser, institutionID, userID)
}

func (uc *SqlUseCase) InsertFacultyUseCase(faculty entity.InstitutionFaculty, institutionID string) error {
	return uc.SqlRepo.InsertFaculty(faculty, institutionID)
}
