package entity

import "project/package/enum"

type PBFTMessage struct {
	VerificationType          enum.VERIFICATIONTYPE     `json:"verification_type"` // PREPREPARE, PREPARE, COMMIT
	OperationID               int                       `json:"operation_id"`
	NodeID                    int                       `json:"node_id"`
	Digest                    []byte                    `json:"digest"`
	Signature                 string                    `json:"signature"`
	Result                    bool                      `json:"result"`
	QRVerificationRequestData QRVerificationRequestData `json:"qr_verification_request_data"`
}

type QRVerificationRequestData struct {
	CertificateHash []byte `json:"certificate_hash"`
	CertificateID   string `json:"certificate_id"` //? certificate id or whwt ?
}

type PBFTExecutionResultEntity struct {
	Result bool  `json:"result"`
	Er     error `json:"error"`
}

func NewPBFTExecutionResultEntity(result bool, er error) PBFTExecutionResultEntity {
	return PBFTExecutionResultEntity{
		Result: result,
		Er:     er,
	}
}
