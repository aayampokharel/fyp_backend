package entity

import "time"

type CertificateData struct {
	// Core Certificate Identity
	CertificateID string `json:"certificate_id"`
	BlockNumber   int    `json:"block_number"`
	Position      int    `json:"position"` // 1-4

	// Student Information (Required)
	StudentID   string `json:"student_id"`
	StudentName string `json:"student_name"`

	// Institution & Faculty Information
	InstitutionID        string `json:"institution_id"`
	InstitutionFacultyID string `json:"institution_faculty_id"`
	PDFCategoryID        string `json:"pdf_category_id"`

	// Certificate Type
	CertificateType string `json:"certificate_type"` // COURSE_COMPLETION, CHARACTER, LEAVING, TRANSFER, PROVISIONAL

	// Academic Information (Optional)
	Degree         string  `json:"degree"`
	College        string  `json:"college"`
	Major          string  `json:"major"`
	GPA            string  `json:"gpa"`
	Percentage     float64 `json:"percentage"`
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
	DataHash        string `json:"data_hash"`
	IssuerPublicKey string `json:"issuer_public_key"`
	CertificateHash string `json:"certificate_hash"` // NEW: Individual certificate hash

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
}

type CertificateDataWithQRCode struct {
	CertificateData `json:"certificate_data"`
	QRCodeBase64    string `json:"qr_code"`
}
