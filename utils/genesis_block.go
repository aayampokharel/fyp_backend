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
			MerkleRoot:   "",
			PreviousHash: "",
			Nonce:        "", //from Pow
			CurrentHash:  "", //from pow
		},
		CertificateData: [4]models.CertificateData{
			{
				ID:                 "admin",
				StudentName:        "admin",
				UniversityName:     "admin",
				Degree:             "admin",
				College:            "admin",
				CertificateDate:    time.Now(),
				Division:           "admin",
				PrincipalSignature: "admin",
				TuApproval:         "admin",
			},
			{},
			{},
			{},
		},
	}
}
