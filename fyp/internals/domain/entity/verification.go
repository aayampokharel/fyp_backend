package entity

import "project/package/enum"

type PBFTMessage struct {
	VerificationType          enum.VERIFICATIONTYPE     `json:"verification_type"` // PREPREPARE, PREPARE, COMMIT
	ViewNumber                int                       `json:"view_number"`
	SequenceNumber            int                       `json:"sequence_number"`
	NodeID                    int                       `json:"node_id"`
	Digest                    string                    `json:"digest"`
	Signature                 string                    `json:"signature"`                    // Digital signature
	QRVerificationRequestData QRVerificationRequestData `json:"qr_verification_request_data"` // Original request
}

type QRVerificationRequestData struct {
	CertificateHash string `json:"certificate_hash"`
	ClientID        string `json:"client_id"` // unique generated ID
	Timestamp       int64  `json:"timestamp"`
}
