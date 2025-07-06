package common

import (
	"github.com/aayampokharel/fyp/models"
	"github.com/aayampokharel/fyp/utils"
)

func CalculateMerkelRoot(allCertificateDatas [4]models.CertificateData) (string, error) {
	length := CalculateLengthOfCertificateArray(allCertificateDatas)
	var allCertificateData []models.CertificateData
	allCertificateData = allCertificateDatas[:]
	if length%2 != 0 {
		allCertificateData = append(allCertificateData, allCertificateDatas[length-1]) //creating copy .
		length++
	}

	var internalHashRoots []string

	for _, val := range allCertificateData {

		hashValue, err := utils.CalculateHashHex(val)
		if err != nil {
			utils.LogErrorWithContext("Hashing", err)
			return "", err
		}
		internalHashRoots = append(internalHashRoots, hashValue)

	}

	for length > 1 {
		var newTempSlice []string
		for key, hashVal := range internalHashRoots {
			if key%2 == 0 {
				newHashValue, err := utils.CalculateHashHex(hashVal + internalHashRoots[key+1])

				if err != nil {
					utils.LogErrorWithContext("HashingLoop", err)
					return "", err
				}
				newTempSlice = append(newTempSlice, newHashValue)
			}

		}
		internalHashRoots = newTempSlice
		length = len(internalHashRoots)
		if length%2 != 0 && length != 1 {
			internalHashRoots = append(internalHashRoots, internalHashRoots[length-1]) //
			length++
		}

	}

	return internalHashRoots[0], nil

}
