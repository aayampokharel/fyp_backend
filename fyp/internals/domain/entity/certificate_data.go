package entity

import "time"

type CertificateData struct {
	// Core Certificate Identity
	CertificateID  string  `json:"certificate_id" gorm:"primaryKey"` // Unique hash of certificate
	StudentID      string  `json:"student_id"`                       // University student ID
	StudentName    string  `json:"student_name"`
	UniversityName string  `json:"university_name"`
	Degree         string  `json:"degree"`
	College        string  `json:"college"`
	Major          string  `json:"major"` // ADDED: Specific field of study
	GPA            string  `json:"gpa"`   // ADDED: Academic performance
	Percentage     float64 `json:"percentage"`
	Division       string  `json:"division"`

	// Dates
	IssueDate      time.Time `json:"issue_date"`      // When cert was issued
	EnrollmentDate time.Time `json:"enrollment_date"` // ADDED: When student started
	CompletionDate time.Time `json:"completion_date"` // ADDED: When course completed

	// Digital Signatures & Verification
	PrincipalSignature string `json:"principal_signature"` // Principal's digital signature

	// Cryptographic Verification
	DataHash        string `json:"data_hash"`         // ADDED: Hash of certificate data
	IssuerPublicKey string `json:"issuer_public_key"` // ADDED: Who issued this certificate

	// Metadata
	CertificateType string    `json:"certificate_type"` // ADDED: "DEGREE", "DIPLOMA", "TRANSCRIPT"
	CreatedAt       time.Time `json:"created_at"`       // ADDED: When this record was created
}
