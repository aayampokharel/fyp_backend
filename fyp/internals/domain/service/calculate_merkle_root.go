package service

import (
	"project/internals/domain/entity"
	"project/package/utils/common"
)

type Service struct {
}

func NewService() Service {
	return Service{}
}

func (s *Service) CalculateMerkleRoot(certificateDataArray [4]entity.CertificateData) (string, error) {
	calculatedLength, err := common.CalculateCertificateDataLength(certificateDataArray)
	if err != nil {
		return "", err
	}
	if calculatedLength == 1 {
		return common.HashData(certificateDataArray[0])
	}

	certificateDataSlice := certificateDataArray[:calculatedLength]
	var hashedCertificateDataSlice []string

	for _, val := range certificateDataSlice {
		hashedVal, err := common.HashData(val)
		if err != nil {
			return "", err
		}
		hashedCertificateDataSlice = append(hashedCertificateDataSlice, hashedVal)

	}
	if calculatedLength%2 != 0 {
		hashedCertificateDataSlice = append(hashedCertificateDataSlice, hashedCertificateDataSlice[calculatedLength-1])
	}
	newLevel := []string{}
	for len(hashedCertificateDataSlice) != 1 {
		if newLevel != nil {
			newLevel = newLevel[:0]

		}

		for i := 1; i <= len(hashedCertificateDataSlice)-1; i += 2 {
			hashedCertificateDataSlice[i], err = common.HashData(hashedCertificateDataSlice[i] + hashedCertificateDataSlice[i-1])
			if err != nil {
				return "", err
			}
			newLevel = append(newLevel, hashedCertificateDataSlice[i])
		}
		hashedCertificateDataSlice = newLevel
	}
	return newLevel[0], nil
}
