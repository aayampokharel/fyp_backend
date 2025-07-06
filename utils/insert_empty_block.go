package utils

import (
	"time"

	"github.com/aayampokharel/fyp/models"
)

func InsertEmptyBlock() models.Block {
	return models.Block{
		Header: models.Header{
			BlockNumber:  0,
			TimeStamp:    time.Now(),
			MerkleRoot:   "",
			PreviousHash: "",
			Nonce:        "", //from Pow
			CurrentHash:  "", //from pow
			Status:       models.PENDING,
		},
		CertificateData: [4]models.CertificateData{
			{},
			{},
			{},
			{},
		}}
}
