package authentication

import (
	"project/internals/domain/entity"
	"project/package/utils/common"
	"strconv"
)

type CreateInstutionResponse struct {
	InstitutionID string `json:"institution_id"`
}

type CreateInstitutionRequest struct {
	InstitutionName string `json:"institution_name"`
	ToleAddress     string `json:"tole_address"`
	WardNumber      int    `json:"ward_number"`
	DistrictAddress string `json:"district_address"`
}

func (c *CreateInstitutionRequest) ToEntity() entity.Institution {
	return entity.Institution{
		InstitutionID:   common.GenerateUUID(16),
		InstitutionName: c.InstitutionName,
		ToleAddress:     c.ToleAddress,
		WardNumber:      strconv.Itoa(c.WardNumber),
		DistrictAddress: c.DistrictAddress,
		IsActive:        false, //true only after inspection from admin
	}
}
