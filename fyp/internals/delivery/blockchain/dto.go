package delivery

import "project/internals/domain/entity"

type CreateCertificateResponse struct {
	Message string `json:"message"`
}
type CreateCertificateDataRequest struct {
	CertificateData []entity.CertificateData `json:"certificate_data"`
}
