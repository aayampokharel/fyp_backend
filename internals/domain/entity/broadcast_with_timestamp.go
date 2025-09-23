package entity

import "time"

type CertificateDataWithTimestamp struct {
	CertificateData CertificateData `json:"certificate_data"`
	Timestamp       time.Time       `json:"timestamp"`
}
