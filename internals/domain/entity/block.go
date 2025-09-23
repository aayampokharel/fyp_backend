package entity

import (
	"project/package/enum"
	"time"
)

type Block struct {
	Header          Header             `json:"header"`
	CertificateData [4]CertificateData `json:"certificate_data"`
}

type Header struct {
	BlockNumber  int         `json:"block_number"`
	TimeStamp    time.Time   `json:"timestamp"`
	PreviousHash string      `json:"previous_hash"`
	Nonce        string      `json:"nonce"`
	CurrentHash  string      `json:"current_hash"`
	MerkleRoot   string      `json:"merkle_root"`
	Status       enum.STATUS `json:"status"`
}
