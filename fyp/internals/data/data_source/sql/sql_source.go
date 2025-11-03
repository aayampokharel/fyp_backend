package sql_source

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"project/internals/data/models"
	"project/internals/domain/entity"
	"project/internals/domain/repository"
	"project/package/enum"
	err "project/package/errors"
	errorz "project/package/errors"
	logger "project/package/utils/pkg"
	"time"

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

func (s *SQLSource) GetPendingInstitutionFromInstitutionID(institutionID string) (*entity.Institution, error) {
	query := `select institution_id, institution_name, tole_address, district_address,ward_number from institutions where institution_id=$1 AND is_active IS NULL AND is_signup_completed=true;`
	var institution entity.Institution
	er := s.DB.QueryRow(query, institutionID).Scan(&institution.InstitutionID, &institution.InstitutionName, &institution.ToleAddress, &institution.DistrictAddress, &institution.WardNumber)
	if er != nil {
		s.logger.Errorln("[sql_source] Error: GetInstitutionFromInstitutionID::", er)
		return nil, er
	}
	return &institution, nil

}
func (s *SQLSource) GetAllPendingInstitutions() ([]entity.Institution, error) {

	query := `select institution_id, institution_name, tole_address, district_address,ward_number from institutions where is_active IS NULL and is_signup_completed=true;`
	rows, err := s.DB.Query(query)
	if err != nil {
		s.logger.Errorln("[sql_source] Error: GetToBeVerifiedInstitutions::", err)
		return nil, err
	}
	defer rows.Close()

	var institutions []entity.Institution
	for rows.Next() {
		var institution entity.Institution
		if err := rows.Scan(&institution.InstitutionID, &institution.InstitutionName, &institution.ToleAddress, &institution.DistrictAddress, &institution.WardNumber); err != nil {
			s.logger.Errorln("[sql_source] Error: GetToBeVerifiedInstitutions::", err)
			return nil, err
		}
		institution.IsActive = false
		institutions = append(institutions, institution)
	}
	return institutions, nil
}

func (s *SQLSource) UpdateSignUpCompletedByInstitutionID(institutionID string) error {
	query := `update institutions set is_signup_completed=true where institution_id=$1;`
	_, err := s.DB.Exec(query, institutionID)
	if err != nil {
		s.logger.Errorln("[sql_source] Error: UpdateFormSubmittedByInstitutionID::", err)
		return err
	}
	return nil
}
func (s *SQLSource) UpdateIsActiveByInstitutionID(institutionID string, isActive bool) error {
	query := `update institutions set is_active=$1 where institution_id=$2;`
	_, err := s.DB.Exec(query, isActive, institutionID)
	if err != nil {
		s.logger.Errorln("[sql_source] Error: UpdateFormSubmittedByInstitutionID::", err)
		return err
	}
	return nil
}
func (s *SQLSource) GetUserIDByInstitutionID(institutionID string) (string, error) {
	var userID string

	query := `
        SELECT user_id 
        FROM institution_user 
        WHERE institution_id = $1;
    `
	err := s.DB.QueryRow(query, institutionID).Scan(&userID)
	if err != nil {
		s.logger.Errorln("[sql_source] Error: GetUserIDByInstitutionID::", err)
		return "", err
	}

	return userID, nil
}

func (s *SQLSource) CheckDuplicationByInstitutionInfo(institution entity.Institution) error {
	var count int
	if institution.InstitutionName == "" || institution.ToleAddress == "" || institution.DistrictAddress == "" || institution.WardNumber == "" {
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
		SELECT id, institution_role, created_at, email, password 
		FROM user_accounts 
		WHERE email = $1 AND deleted_at IS NULL
		LIMIT 1
	`
	var userAccount entity.UserAccount
	er := s.DB.QueryRow(query, userEmail).Scan(
		&userAccount.ID,
		&userAccount.InstitutionRole,
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

	query := `INSERT INTO institutions (institution_id, institution_name, tole_address, district_address,ward_number)
		VALUES ($1, $2, $3, $4, $5);`

	if _, er := s.DB.Exec(query, institution.InstitutionID, institution.InstitutionName, institution.ToleAddress, institution.DistrictAddress, institution.WardNumber); er != nil {
		s.logger.Errorln("[sql_source] Error: InsertInstitutions::", er)
		return "", er
	}
	return institution.InstitutionID, nil

}

func (s *SQLSource) InsertUserAccounts(userAccounts entity.UserAccount) (string, error) {
	userAccountModel := models.UserAccountFromEntity(userAccounts)
	var createdAt string
	insertQuery := `Insert into user_accounts (id,system_role, institution_role, email, password) values ($1, $2, $3, $4,$5)  RETURNING created_at;`
	if er := s.DB.QueryRow(insertQuery, userAccountModel.ID, userAccountModel.SystemRole, userAccountModel.InstitutionRole, userAccountModel.Email, userAccountModel.Password).Scan(&createdAt); er != nil {
		s.logger.Errorln("[sql_source] Error: InsertUserAccounts::", er)
		return "", er
	}
	return createdAt, nil

}

func (s *SQLSource) InsertInstitutionUser(institutionUser entity.InstitutionUser) error {
	query := `INSERT INTO institution_user (institution_id, user_id, institution_logo_base64) VALUES ($1, $2, $3);`

	_, err := s.DB.Exec(query, institutionUser.InstitutionID, institutionUser.UserID, institutionUser.InstitutionLogoBase64)
	if err != nil {
		s.logger.Errorln("[sql_source] Error: InsertInstitutionUser::", query+"::"+institutionUser.InstitutionID+"::"+institutionUser.UserID+"::"+institutionUser.InstitutionLogoBase64)
		s.logger.Errorln("[sql_source] Error: InsertInstitutionUser::", err)
		return err
	}
	return nil
}

func (s *SQLSource) VerifyAdminLogin(userMail, password string) (string, time.Time, error) {
	////CHECK FOR USER_NOT FOUND RETURN ERROR ,errorz package used here
	//// incorrect password case as well for matching email.
	//! or simply user not found

	query := `SELECT id,created_at FROM user_accounts WHERE email=$1 AND password=$2 AND system_role=$3 AND institution_role IS NULL AND deleted_at IS NULL;`
	var adminID string
	var createdAt time.Time
	err := s.DB.QueryRow(query, userMail, password, enum.ADMIN.String()).Scan(&adminID, &createdAt)
	if err != nil {
		s.logger.Errorln("[sql_source] Error: VerifyAdminLogin::", err)
		return "", time.Time{}, err
	}
	return adminID, createdAt, nil

}
func (s *SQLSource) InsertFaculty(faculty entity.InstitutionFaculty) (facultyID string, er error) {
	query := `
		INSERT INTO institution_faculty (
			institution_faculty_id,
			institution_id,
			faculty_name,
			faculty_authority_with_signature,
			faculty_public_key,
			university_affiliation,
			university_college_code
		) VALUES ($1, $2, $3, $4, $5, $6, $7);
	`
	signatureJSON, er := json.Marshal(faculty.FacultyAuthorityWithSignatures)
	if er != nil {
		return "", er
	}
	_, er = s.DB.Exec(
		query,
		faculty.InstitutionFacultyID,
		faculty.InstitutionID,
		faculty.FacultyName,
		signatureJSON,
		faculty.FacultyPublicKey,
		faculty.UniversityAffiliation,
		faculty.UniversityCollegeCode,
	)
	if er != nil {
		s.logger.Errorln("[sql_source] Error: InsertFaculty::", er)
		return "", er
	}
	return faculty.InstitutionFacultyID, nil
}

func (s *SQLSource) InsertPDFFile(pdfFile entity.PDFFileEntity) error {
	query := `INSERT INTO pdf_files (file_id,category_id, file_name, pdf_data) VALUES ($1, $2, $3,$4);`

	_, err := s.DB.Exec(query, pdfFile.FileID, pdfFile.CategoryID, pdfFile.FileName, pdfFile.PDFData)
	if err != nil {
		s.logger.Errorln("[sql_source] Error: InsertPDFFile::", err)
		return err
	}
	return nil
}

func (s *SQLSource) InsertAndGetPDFCategory(pdfFileCategory entity.PDFFileCategoryEntity) (*entity.PDFFileCategoryEntity, error) {
	query := `
		INSERT INTO pdf_file_categories 
		(category_id, category_name, institution_id, institution_faculty_id)
		VALUES ($1, $2, $3, $4)
		RETURNING category_id, category_name, institution_id, institution_faculty_id;
	`

	var inserted entity.PDFFileCategoryEntity
	err := s.DB.QueryRow(
		query,
		pdfFileCategory.CategoryID,
		pdfFileCategory.CategoryName,
		pdfFileCategory.InstitutionID,
		pdfFileCategory.InstitutionFacultyID,
	).Scan(
		&inserted.CategoryID,
		&inserted.CategoryName,
		&inserted.InstitutionID,
		&inserted.InstitutionFacultyID,
	)

	if err != nil {
		s.logger.Errorln("[sql_source] Error: InsertAndGetPDFCategory::", err)
		return nil, err
	}

	return &inserted, nil
}

func (s *SQLSource) RetrievePDFFileByFileIDOrCategoryID(fileID string, categoryID string, isDownloadAll bool) ([]entity.PDFFileEntity, error) {

	var (
		pdfFileEntities     []entity.PDFFileEntity
		singlePdfFileEntity entity.PDFFileEntity
		query               string
		rows                *sql.Rows
		er                  error
	)
	if !isDownloadAll {
		query = `select file_id,pdf_data,file_name,uploaded_at from pdf_files where file_id=$1 AND category_id=$2;`
		rows, er = s.DB.Query(query, fileID, categoryID)
		if er != nil {
			s.logger.Errorln("[sql_source] Error: retrievePDFfileByFileIDOrCategoryID", er)
			return nil, er
		}

	} else {
		query = `select file_id,pdf_data,file_name,uploaded_at from pdf_files where category_id=$1;`
		rows, er = s.DB.Query(query, categoryID)
		if er != nil {
			s.logger.Errorln("[sql_source] Error: retrievePDFfileByFileIDOrCategoryID", er)
			return nil, er
		}
	}
	for rows.Next() {
		er = rows.Scan(&singlePdfFileEntity.FileID, &singlePdfFileEntity.PDFData, &singlePdfFileEntity.FileName, &singlePdfFileEntity.UploadedAt)
		if er != nil {
			s.logger.Errorln("[sql_source] Error: retrievePDFfileByFileIDOrCategoryID", er)
			return nil, er
		}
		pdfFileEntities = append(pdfFileEntities, singlePdfFileEntity)
	}
	if er = rows.Err(); er != nil {
		s.logger.Errorln("[sql_source] Error during rows iteration:", er)
		return nil, er
	}
	return pdfFileEntities, nil
}

func (s *SQLSource) InsertBlockWithSingleCertificate(blockHeader entity.Header, certificateData entity.CertificateData, certificatePositionZeroIndex int) error {

	if certificatePositionZeroIndex < 0 || certificatePositionZeroIndex >= 4 {
		s.logger.Errorln("[sql_source] certificate position index error::", certificatePositionZeroIndex)
		return err.ErrArrayOutOfBound

	}
	// Start transaction
	tx, err := s.DB.Begin()
	if err != nil {
		s.logger.Errorln("[sql_source] Error starting transaction:", err)
		return err
	}
	defer tx.Rollback() // Safe error rollback

	blockQuery := `
        INSERT INTO blocks (
            block_number, timestamp, previous_hash, nonce, 
            current_hash, merkle_root, status, created_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `

	_, err = tx.Exec(
		blockQuery,
		blockHeader.BlockNumber,
		blockHeader.TimeStamp,
		blockHeader.PreviousHash,
		blockHeader.Nonce,
		blockHeader.CurrentHash,
		blockHeader.MerkleRoot,
		time.Now(),
	)
	if err != nil {
		s.logger.Errorln("[sql_source] Error inserting block:", err)
		return err
	}

	er := s.InsertCertificate(certificateData, blockHeader.BlockNumber, certificatePositionZeroIndex)
	if er != nil {
		s.logger.Errorln("[sql_source] Error inserting certificate:", er)
		return er
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		s.logger.Errorln("[sql_source] Error committing transaction:", err)
		return err
	}

	s.logger.Infof("[sql_source] Successfully inserted block %d with certificate", blockHeader.BlockNumber)
	return nil
}

func (s *SQLSource) InsertCertificate(certificateData entity.CertificateData, blockNumber int, certificatePositionZeroIndex int) error {

	certQuery := `
        INSERT INTO certificates (
            certificate_id, block_number, position,
            student_id, student_name, 
            institution_id, institution_faculty_id, pdf_category_id,
            certificate_type,
            degree, college, major, gpa, percentage, division, university_name,
            issue_date, enrollment_date, completion_date, leaving_date,
            reason_for_leaving, character_remarks, general_remarks,
            data_hash, faculty_public_key, created_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26)
    `

	_, er := s.DB.Exec(
		certQuery,
		certificateData.CertificateID,
		blockNumber,
		certificatePositionZeroIndex+1, // Position 1-4

		// Student Information
		certificateData.StudentID,
		certificateData.StudentName,

		// Institution Information
		certificateData.InstitutionID,
		certificateData.InstitutionFacultyID,
		certificateData.PDFCategoryID,

		// Certificate Type
		certificateData.CertificateType,

		// Academic Information
		certificateData.Degree,
		certificateData.College,
		certificateData.Major,
		certificateData.GPA,
		certificateData.Percentage,
		certificateData.Division,
		certificateData.UniversityName,

		// Date Information
		certificateData.IssueDate,
		certificateData.EnrollmentDate,
		certificateData.CompletionDate,
		certificateData.LeavingDate,

		// Reason Fields
		certificateData.ReasonForLeaving,
		certificateData.CharacterRemarks,
		certificateData.GeneralRemarks,

		// Cryptographic Verification
		certificateData.CertificateHash,
		certificateData.FacultyPublicKey, // Changed from IssuerPublicKey

		// Timestamp
		certificateData.CreatedAt,
	)
	if er != nil {
		s.logger.Errorln("[sql_source] Error inserting certificate at position", certificatePositionZeroIndex+1, ":", er)
		return er
	}

	s.logger.Infof("[sql_source] Successfully inserted certificate for block %d ", blockNumber)
	return nil
}

func (s *SQLSource) UpdateBlockDataByID(blockHeader entity.Header, id string) error {
	query := `update blocks set timestamp=$1, previous_hash=$2, nonce=$3, current_hash=$4, merkle_root=$5 where block_number=$7;`
	_, err := s.DB.Exec(query, blockHeader.TimeStamp, blockHeader.PreviousHash, blockHeader.Nonce, blockHeader.CurrentHash, blockHeader.MerkleRoot, id)
	if err != nil {
		s.logger.Errorln("[sql_source] Error: UpdateBlockDataByID::", err)
		return err
	}
	return nil
}

func (s *SQLSource) GetFacultyPublicKey(id string) (string, error) {
	query := `select faculty_public_key from institution_faculty where institution_faculty_id=$1;`
	var facultyPublicKey string
	if er := s.DB.QueryRow(query, id).Scan(&facultyPublicKey); er != nil {
		s.logger.Errorln("[sql_source] Error: GetIssuerPublicKey::", er)
		return "", er
	}
	return facultyPublicKey, nil

}
func (s *SQLSource) GetInfoFromPdfFilesCategories(categoryID string) (*entity.PDFFileCategoryEntity, error) {
	query := `select * from pdf_file_categories where category_id=$1;`
	var pdfFileCategory entity.PDFFileCategoryEntity
	if er := s.DB.QueryRow(query, categoryID).Scan(&pdfFileCategory.CategoryID, &pdfFileCategory.CategoryName, &pdfFileCategory.InstitutionID, &pdfFileCategory.InstitutionFacultyID, &pdfFileCategory.CreatedAt); er != nil {
		s.logger.Errorln("[sql_source] Error: GetInfoFromPdfFilesCategories::", er)
		return nil, er
	}
	return &pdfFileCategory, nil
}
