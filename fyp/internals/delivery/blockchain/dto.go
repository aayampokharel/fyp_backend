package delivery

import (
	"encoding/hex"
	"project/internals/domain/entity"
	err "project/package/errors"
	"project/package/utils/common"
	"time"
)

// type CreateIndividualCertificateResponse struct {
// 	Message       string `json:"message"`
// 	FileID        string `json:"file_id"`
// 	StudentName   string `json:"student_name"`
// 	FileName      string `json:"file_name"`
// StudentID     string `json:"student_id"`
// CertificateID string `json:"certificate_id"`
// CategoryID    string `json:"category_id"`

// }

type BasicStudentInfoDto struct {
	StudentID   string `json:"student_id"`
	StudentName string `json:"student_name"`
	FileID      string `json:"file_id"`
	FileName    string `json:"file_name"`
	FacultyName string `json:"faculty_name"`
	////remainiing later .
}
type CreateAllCertificateResponse struct {
	Message string `json:"message"`
	// StudentList []BasicStudentInfoDto `json:"student_list"`
}

type MinimalCertificateData struct {

	//BlockNumber   int    `json:"block_number"`
	//Position      int    `json:"position"` // 1-4

	// Student Information (Required)
	StudentID   string `json:"student_id"`
	StudentName string `json:"student_name"`

	// Institution & Faculty Information
	InstitutionID        string `json:"institution_id"`
	InstitutionFacultyID string `json:"institution_faculty_id"`
	// //PDFCategoryID        string `json:"pdf_category_id"`  "this isnot made at the time . i have to create it before doing anything."

	// Certificate Type
	CertificateType string `json:"certificate_type"` // COURSE_COMPLETION, CHARACTER, LEAVING, TRANSFER, PROVISIONAL

	// Academic Information (Optional)
	Degree         string  `json:"degree"`
	College        string  `json:"college"`
	Major          string  `json:"major"`
	GPA            string  `json:"gpa"`
	Percentage     *string `json:"percentage"`
	Division       string  `json:"division"`
	UniversityName string  `json:"university_name"`

	// Date Information
	IssueDate      time.Time `json:"issue_date"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	CompletionDate time.Time `json:"completion_date"`
	LeavingDate    time.Time `json:"leaving_date"`

	// Reason Fields
	ReasonForLeaving string `json:"reason_for_leaving"`
	CharacterRemarks string `json:"character_remarks"`
	GeneralRemarks   string `json:"general_remarks"`

	// Cryptographic Verification
	//DataHash        string `json:"data_hash"`
	//IssuerPublicKey string `json:"issuer_public_key"`
	//CertificateHash string `json:"certificate_hash"` // NEW: Individual certificate hash

	// Timestamps
	//CreatedAt time.Time `json:"created_at"`
}

type GetRandomDataInsertionRequest []string

const RandomID = "random_id"

var GetAllPendingInstitutionsQuery = GetRandomDataInsertionRequest{RandomID}

type CreateCertificateDataRequest struct {
	InstitutionID          string                   `json:"institution_id"`
	InstitutionFacultyID   string                   `json:"institution_faculty_id"`
	InstitutionFacultyName string                   `json:"institution_faculty_name"`
	CategoryID             string                   `json:"category_id"`
	CategoryName           string                   `json:"category_name"`
	CertificateData        []MinimalCertificateData `json:"certificate_data"`
}

func (m *CreateCertificateDataRequest) ToPdfFileCategoryEntity() (entity.PDFFileCategoryEntity, error) {

	if m.InstitutionFacultyID == "" || m.InstitutionID == "" || m.CategoryName == "" {
		return entity.PDFFileCategoryEntity{}, err.ErrEmptyFields
	}
	return entity.PDFFileCategoryEntity{
		CategoryID:           common.GenerateUUID(16),
		CategoryName:         m.CategoryName,
		InstitutionID:        m.InstitutionID,
		InstitutionFacultyID: m.InstitutionFacultyID,
	}, nil
}

func (m *MinimalCertificateData) ToEntity(categoryID string, institutionName, universityAffiliation string) (*entity.CertificateData, error) {

	var percentageFloat float64
	var er error
	if m.Percentage != nil {
		percentageFloat, er = common.ConvertToFloat(*m.Percentage)
		if er != nil {

			return nil, er
		}

	}
	certificateData := entity.CertificateData{
		// BlockNumber:          blockNumber, //being done inside usecase
		CertificateID:        common.GenerateUUID(16),
		PDFFileID:            common.GenerateUUID(16),
		PDFCategoryID:        categoryID,
		StudentID:            m.StudentID,
		StudentName:          m.StudentName,
		InstitutionID:        m.InstitutionID,
		InstitutionFacultyID: m.InstitutionFacultyID,
		CertificateType:      m.CertificateType, // COURSE_COMPLETION, CHARACTER, LEAVING, TRANSFER, PROVISIONAL

		// Academic Information (Optional)

		Degree:         m.Degree,
		College:        institutionName,
		Major:          m.Major,
		GPA:            m.GPA,
		Percentage:     &percentageFloat,
		Division:       m.Division,
		UniversityName: universityAffiliation,

		// Date Information
		IssueDate:      m.IssueDate,
		EnrollmentDate: m.EnrollmentDate,
		CompletionDate: m.CompletionDate,
		LeavingDate:    m.LeavingDate,

		// Reason Fields
		ReasonForLeaving: m.ReasonForLeaving,
		CharacterRemarks: m.CharacterRemarks,
		GeneralRemarks:   m.GeneralRemarks,
		CreatedAt:        time.Now(),

		// Cryptographic Verification
		//DataHash        string `json:"data_hash"`
		//IssuerPublicKey string `json:"issuer_public_key"`
		//CertificateHash string `json:"certificate_hash"` // NEW: Individual certificate hash

		// Timestamps
		//CreatedAt time.Time `json:"created_at"`
	}
	dataToHashString := certificateData.GetCertificateDataForHash()
	_, hash, er := common.HashData(dataToHashString)
	if er != nil {

		return nil, er
	}

	certificateData.CertificateHash = hex.EncodeToString(hash)
	return &certificateData, nil
}

func FromPDFFileCategoryToPDFFileEntity(categoryID string, studentName, faculty string, fileID string, index int) entity.PDFFileEntity {
	return entity.PDFFileEntity{
		CategoryID: categoryID,
		FileID:     fileID,
		FileName:   common.GeneratePDFFileName(studentName, faculty, index),
	}
}

type GetCertificateDataListRequest []string

const InstitutionID string = "institution_id"
const InstitutionFacultyID string = "institution_faculty_id"
const CategoryID string = "category_id"

var GetCertificateDataListRequestQuery = GetCertificateDataListRequest{InstitutionID, InstitutionFacultyID, CategoryID}
