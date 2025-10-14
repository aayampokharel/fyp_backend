package repository

import "project/internals/domain/entity"

type ISqlRepository interface {
	GetUserAccountsForEmail(userEmail string) (entity.UserAccount, error)
	InsertInstitutions(institution entity.Institution) error
	InsertUserAccounts(userAccounts entity.UserAccount) error
	InsertInstitutionUser(institutionUser entity.InstitutionUser, institutionID, userID string) error
	InsertFaculty(faculty entity.InstitutionFaculty, institutionID string) error
}
