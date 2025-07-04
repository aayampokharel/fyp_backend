package utils

import (
	"github.com/aayampokharel/fyp/models"
)

func CalculateMerkelRoot(allCertificateData []models.CertificateData) (string, error) {
	length := len(allCertificateData)
	if length%2 != 0 {
		allCertificateData = append(allCertificateData, allCertificateData[length-1]) //creating copy .
		length++
	}

	var internalHashRoots []string

	for _, val := range allCertificateData {
		hashValue, err := CalculateHashHex(val)
		if err != nil {
			return "", err
		}
		internalHashRoots = append(internalHashRoots, hashValue)

	}
	for length > 1 {
		var newTempSlice []string
		for key, hashVal := range internalHashRoots {
			if key%2 == 0 {
				newHashValue, err := CalculateHashHex(hashVal + internalHashRoots[key+1])

				if err != nil {
					return "", err
				}
				newTempSlice = append(newTempSlice, newHashValue)
			}

		}
		internalHashRoots = newTempSlice
		length = len(internalHashRoots)
		if length%2 != 0 {
			internalHashRoots = append(internalHashRoots, internalHashRoots[length-1]) //creating copy .
			length++
		}
	}
	return internalHashRoots[0], nil

}
