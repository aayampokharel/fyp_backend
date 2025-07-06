package common

import "github.com/aayampokharel/fyp/models"

func CalculateLengthOfCertificateArray(array [4]models.CertificateData) int {
	emptyCertificateData := models.CertificateData{}
	count := 0
	for _, value := range array {
		if value != emptyCertificateData {
			count++
		}
	}
	return count
}
