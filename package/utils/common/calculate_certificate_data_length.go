package common

import (
	"project/internals/domain/entity"
	err "project/package/errors"
)

func CalculateCertificateDataLength(certificateDataArray [4]entity.CertificateData) (int, error) {
	certificateDataSlice := certificateDataArray[:]
	count := 0
	emptyValue := entity.CertificateData{}

	for _, value := range certificateDataSlice {
		if value != emptyValue {
			count++
		}
	}

	if count < 0 || count > len(certificateDataArray)-1 {
		return -1, err.ErrArrayOutOfBound
	}
	return count, nil
}
