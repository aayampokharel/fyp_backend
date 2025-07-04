package utils

import (
	"time"

	"github.com/aayampokharel/fyp/models"
)

func CreateGenesisBlock() models.Block {
	return models.Block{
		Header: models.Header{
			BlockNumber:  0,
			TimeStamp:    time.Now(),
			PreviousHash: "",
			Nonce:        "",
			CurrentHash:  "",
			MerkleRoot:   "",
		},
		CertificateData: [4]models.CertificateData{
			{
				ID:                 "",
				StudentName:        "",
				UniversityName:     "",
				Degree:             "",
				College:            "",
				CertificateDate:    time.Now(),
				Division:           "",
				PrincipalSignature: "",
				TuApproval:         "",
			},
		},
	}
}
