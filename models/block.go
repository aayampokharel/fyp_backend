package models

import (
	"time"
)

type Block struct {
	Header          Header             `json:"header"`
	CertificateData [4]CertificateData `json:"certificate_data"`
}

type Header struct {
	BlockNumber  uint64    `json:"block_number"`
	TimeStamp    time.Time `json:"timestamp"`
	PreviousHash string    `json:"previous_hash"`
	Nonce        string    `json:"nonce"`
	CurrentHash  string    `json:"current_hash"`
	MerkleRoot   string    `json:"merkle_root"`
	Status       STATUS    `json:"status"`
}
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

type BlockWithSignature struct {
	BlockData Block  `json:"block_data"`
	Signature string `json:"signature"`
}
