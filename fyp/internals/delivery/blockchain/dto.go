package delivery

import "project/internals/domain/entity"

type CreateCertificateResponse struct {
	Body []entity.Block `json:"body"`
}
