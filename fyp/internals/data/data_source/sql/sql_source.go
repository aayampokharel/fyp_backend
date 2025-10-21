package sql_source

import (
	"database/sql"
	"fmt"
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	errorz "project/package/errors"
	logger "project/package/utils/pkg"

	"go.uber.org/zap"
)

type SQLSource struct {
	DB     *sql.DB
	logger *zap.SugaredLogger
}

func NewSQLSource(db *sql.DB) *SQLSource {
	return &SQLSource{DB: db, logger: logger.Logger}
}

var _ repository.ISqlRepository = (*SQLSource)(nil)

func (s *SQLSource) CheckDuplicationByInstitutionInfo(institution entity.Institution) error {
	var count int
	if institution.InstitutionName == "" || institution.ToleAddress == "" || institution.DistrictAddress == "" {
		s.logger.Errorln("[sql_source] Error: CheckDuplicationByInstitutionInfo::", errorz.ErrEmptyInstitutionInfo)
		return errorz.ErrEmptyInstitutionInfo
	}
	query := `SELECT count(*) FROM institutions WHERE institution_name=$1 and tole_address=$2 and ward_number=$3 and district_address=$4 and is_active=true;`

	err := s.DB.QueryRow(query, institution.InstitutionName, institution.ToleAddress, institution.WardNumber, institution.DistrictAddress).Scan(&count)
	if err != nil {
		s.logger.Errorln("[sql_source] Error: CheckDuplicationByInstitutionInfo::", err)
		return err
	}
	if count > 0 {
		return errorz.ErrInstitutionAlreadyRegistered
	}
	return nil

}
func (s *SQLSource) GetUserAccountsForEmail(userEmail string) (entity.UserAccount, error) {

	if userEmail == "" {
		s.logger.Errorln("[sql_source] Error: GetUserAccountsForEmail::", errorz.ErrEmptyUserEmail)
		return entity.UserAccount{}, errorz.ErrEmptyUserEmail
	}

	query := `
		SELECT id, role, created_at, email, password 
		FROM user_accounts 
		WHERE email = $1 AND deleted_at IS NULL
		LIMIT 1
	`
	var userAccount entity.UserAccount
	er := s.DB.QueryRow(query, userEmail).Scan(
		&userAccount.ID,
		&userAccount.Role,
		&userAccount.CreatedAt,
		&userAccount.Email,
		&userAccount.Password)
	if er != nil {
		if er == sql.ErrNoRows {
			return entity.UserAccount{}, errorz.ErrWithMoreInfo(er, fmt.Sprintf("user with email '%s' not found or is deleted", userEmail))
		}
		s.logger.Errorln("[sql_source] Error: GetUserAccountsForEmail::", er)
		return entity.UserAccount{}, er
	}

	return userAccount, nil
}

func (s *SQLSource) InsertInstitutions(institution entity.Institution) (string, error) {

	query := `INSERT INTO institutions (institution_id, institution_name, tole_address, district_address,ward_number,is_active)
		VALUES ($1, $2, $3, $4, $5,$6);`

	if _, er := s.DB.Exec(query, institution.InstitutionID, institution.InstitutionName, institution.ToleAddress, institution.DistrictAddress, institution.WardNumber, institution.IsActive); er != nil {
		s.logger.Errorln("[sql_source] Error: InsertInstitutions::", er)
		return "", er
	}
	return institution.InstitutionID, nil

}

func (s *SQLSource) InsertUserAccounts(userAccounts entity.UserAccount) error {
	insertQuery := `Insert into user_accounts (id, role, email, password) values ($1, $2, $3, $4);`
	if _, er := s.DB.Exec(insertQuery, userAccounts.ID, userAccounts.Role, userAccounts.Email, userAccounts.Password); er != nil {
		s.logger.Errorln("[sql_source] Error: InsertUserAccounts::", er)
		return er
	}
	return nil

}

func (s *SQLSource) InsertInstitutionUser(institutionUser entity.InstitutionUser) error {
	query := `INSERT INTO institution_user (institution_id, user_id, public_key,  institution_logo_base64) VALUES ($1, $2, $3, $4, $5, $6);`

	_, err := s.DB.Exec(query, institutionUser.InstitutionID, institutionUser.UserID, institutionUser.PublicKey, institutionUser.InstitutionLogoBase64)
	if err != nil {
		s.logger.Errorln("[sql_source] Error: InsertInstitutionUser::", err)
		return err
	}
	return nil

}
func (s *SQLSource) InsertFaculty(faculty entity.InstitutionFaculty, institutionID string) error {
	query := `INSERT INTO institution_faculty (institution_faculty_id, institution_id, faculty, faculty_hod_name, faculty_hod_signature_base64) VALUES ($1, $2, $3, $4, $5);`
	_, err := s.DB.Exec(query, faculty.InstitutionFacultyID, institutionID, faculty.Faculty, faculty.FacultyHODName, faculty.FacultyHODSignatureBase64)
	if err != nil {
		s.logger.Errorln("[sql_source] Error: InsertFaculty::", err)
		return err
	}
	return nil
}
