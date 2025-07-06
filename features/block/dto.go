package block

import (
	"crypto/sha256"
	"encoding/json"

	"github.com/aayampokharel/fyp/models"
)

func ConvertStatusEnumToString(status models.STATUS) string {
	switch status {
	case models.APPROVED:
		return "APPROVED"
	case models.REJECTED:
		return "REJECTED"
	case models.PENDING:
		return "PENDING"
	default:
		return ""
	}
}
func ConvertStringStatusToEnum(status string) models.STATUS {
	switch status {
	case "APPROVED":
		return models.APPROVED
	case "REJECTED":
		return models.REJECTED
	case "PENDING":
		return models.PENDING
	default:
		return ""
	}
}

func HashCertificateData(block models.CertificateData) ([]byte, error) {
	jsonData, err := json.Marshal(block)
	if err != nil {
		return nil, err
	}
	hashedValue := sha256.Sum256(jsonData)
	return hashedValue[:], nil
}

func ConvertintoModelsBlock(blockDataWithSignature models.CertificateDataWithSignature) models.CertificateData {
	return blockDataWithSignature.CertificateData
}
