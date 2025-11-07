package entity

import "time"

type HashableData struct {
	CertificateID        string    `json:"certificate_id"`
	StudentID            string    `json:"student_id"`
	StudentName          string    `json:"student_name"`
	InstitutionID        string    `json:"institution_id"`
	InstitutionFacultyID string    `json:"institution_faculty_id"`
	UniversityName       string    `json:"university_name"`
	Degree               string    `json:"degree"`
	College              string    `json:"college"`
	Major                string    `json:"major"`
	GPA                  string    `json:"gpa"`
	Percentage           float64   `json:"percentage"`
	ReasonForLeaving     string    `json:"reason_for_leaving"`
	CharacterRemarks     string    `json:"character_remarks"`
	GeneralRemarks       string    `json:"general_remarks"`
	Division             string    `json:"division"`
	EnrollmentDate       time.Time `json:"enrollment_date"`
	CompletionDate       time.Time `json:"completion_date"`
	IssueDate            time.Time `json:"issue_date"`
	CertificateType      string    `json:"certificate_type"`
	FacultyPublicKey     string    `json:"faculty_public_key"`
}
