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
	Degree         string  `json:"degree,omitempty"`
	College        string  `json:"college,omitempty"`
	Major          string  `json:"major,omitempty"`
	GPA            string  `json:"gpa,omitempty"`
	Percentage     float64 `json:"percentage,omitempty"`
	Division       string  `json:"division,omitempty"`
	UniversityName string  `json:"university_name,omitempty"`

	// Date Information
	IssueDate      time.Time `json:"issue_date"`
	EnrollmentDate time.Time `json:"enrollment_date,omitempty"`
	CompletionDate time.Time `json:"completion_date,omitempty"`
	LeavingDate    time.Time `json:"leaving_date,omitempty"`

	// Reason Fields
	ReasonForLeaving string `json:"reason_for_leaving,omitempty"`
	CharacterRemarks string `json:"character_remarks,omitempty"`
	GeneralRemarks   string `json:"general_remarks,omitempty"`

	// Cryptographic Verification
	CertificateHash  string `json:"certificate_hash"`
	FacultyPublicKey string `json:"faculty_public_key"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
}

type CertificateDataWithQRCode struct {
	CertificateData `json:"certificate_data"`
	QRCodeBase64    string `json:"qr_code"`
}

func (c *CertificateData) ToHashableData() *HashableData {
	return &HashableData{
		CertificateID:        c.CertificateID,
		StudentID:            c.StudentID,
		StudentName:          c.StudentName,
		InstitutionID:        c.InstitutionID,
		InstitutionFacultyID: c.InstitutionFacultyID,
		UniversityName:       c.UniversityName,
		Degree:               c.Degree,
		College:              c.College,
		Major:                c.Major,
		GPA:                  c.GPA,
		Division:             c.Division,
		EnrollmentDate:       c.EnrollmentDate,
		CompletionDate:       c.CompletionDate,
		IssueDate:            c.IssueDate,
		CertificateType:      c.CertificateType,
		FacultyPublicKey:     c.FacultyPublicKey,
		ReasonForLeaving:     c.ReasonForLeaving,
		CharacterRemarks:     c.CharacterRemarks,
		GeneralRemarks:       c.GeneralRemarks,
	}
}
