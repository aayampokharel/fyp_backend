package models

import (
	"project/internals/domain/entity"
)

type Institution struct {
	InstitutionID   string `json:"institution_id"`
	InstitutionName string `json:"institution_name"`
	ToleAddress     string `json:"tole_address"`
	WardNumber      string `json:"ward_number"`
	DistrictAddress string `json:"district_address"`
	IsActive        *bool  `json:"is_active"`
}

func (i *Institution) ToEntity() entity.Institution {
	return entity.Institution{
		InstitutionID:   i.InstitutionID,
		InstitutionName: i.InstitutionName,
		ToleAddress:     i.ToleAddress,
		WardNumber:      i.WardNumber,
		DistrictAddress: i.DistrictAddress,
		IsActive:        i.IsActive,
	}
}
