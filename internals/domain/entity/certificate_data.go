package entity

import "time"

type CertificateData struct {
	ID                 string    `json:"id"`
	StudentName        string    `json:"student_name"`
	UniversityName     string    `json:"university_name"`
	Degree             string    `json:"degree"`
	College            string    `json:"college"`
	CertificateDate    time.Time `json:"date"`
	Division           string    `json:"division"`
	PrincipalSignature string    `json:"principal_signature"` //principal's digital signature
	TuApproval         string    `json:"tu_approval"`         //tu's digital signature
}
