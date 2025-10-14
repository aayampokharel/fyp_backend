package sql_source

import (
	"database/sql"
	"fmt"
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	err "project/package/errors"
	erroz "project/package/errors"
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

func (s *SQLSource) GetUserAccountsForEmail(userEmail string) (entity.UserAccount, error) {

	if userEmail == "" {
		s.logger.Errorln("[sql_source] Error: GetUserAccountsForEmail::", erroz.ErrEmptyUserEmail)
		return entity.UserAccount{}, erroz.ErrEmptyUserEmail
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
			return entity.UserAccount{}, err.ErrWithMoreInfo(er, fmt.Sprintf("user with email '%s' not found or is deleted", userEmail))
		}
		s.logger.Errorln("[sql_source] Error: GetUserAccountsForEmail::", er)
		return entity.UserAccount{}, er
	}

	return userAccount, nil
}

func (s *SQLSource) InsertInstitutions(institution entity.Institution) error {
	institutionCheckQuery := `
SELECT institution_id FROM institutions WHERE institution_name=$1 and tole_address=$2 and district_address=$3 and is_active=true;`

	var institutionID string
	err := s.DB.QueryRow(institutionCheckQuery, institution.InstitutionName, institution.ToleAddress, institution.DistrictAddress).Scan(&institutionID)

	if err == sql.ErrNoRows {

		query := `INSERT INTO institutions (institution_id, institution_name, tole_address, district_address, is_active)
		VALUES ($1, $2, $3, $4, $5);`

		if _, er := s.DB.Exec(query, institution.InstitutionID, institution.InstitutionName, institution.ToleAddress, institution.DistrictAddress, institution.IsActive); er != nil {
			s.logger.Errorln("[sql_source] Error: InsertInstitutions::", er)
			return er
		}
		return nil
	}

	if err != nil {
		s.logger.Errorln("[sql_source] Error checking institution existence:", err)
		return err
	}

	// Institution already exists
	return erroz.ErrWithMoreInfo(err, fmt.Sprintf("institution already exists with ID: %s", institutionID))
}

func (s *SQLSource) InsertUserAccounts(userAccounts entity.UserAccount) error {
	//check if email registered
	var userId string
	checkEmailQuery := `Select id from user_accounts where email=$1 and deleted_at is null;`

	err := s.DB.QueryRow(checkEmailQuery, userAccounts.Email).Scan(&userId)

	if err == sql.ErrNoRows {
		insertQuery := `Insert into user_accounts (id, role, email, password) values ($1, $2, $3, $4);`
		if _, er := s.DB.Exec(insertQuery, userAccounts.ID, userAccounts.Role, userAccounts.Email, userAccounts.Password); er != nil {
			s.logger.Errorln("[sql_source] Error: InsertUserAccounts::", er)
			return er
		}
		return nil
	}

	if err != nil {
		s.logger.Errorln("[sql_source] Error: InsertUserAccounts::", err)
		return err
	}

	return erroz.ErrWithMoreInfo(err, fmt.Sprintf("user with email '%s' already exists", userAccounts.Email))
}

func (s *SQLSource) InsertInstitutionUser(institutionUser entity.InstitutionUser, institutionID, userID string) error {
	query := `INSERT INTO institution_user (institution_id, user_id, public_key, principal_name, principal_signature_base64, institution_logo_base64) VALUES ($1, $2, $3, $4, $5, $6);`

	_, err := s.DB.Exec(query, institutionID, userID, institutionUser.PublicKey, institutionUser.PrincipalName, institutionUser.PrincipalSignatureBase64, institutionUser.InstitutionLogoBase64)
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
