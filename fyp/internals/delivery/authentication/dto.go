package authentication

import (
	"project/internals/domain/entity"
	"project/package/enum"
	"project/package/utils/common"
	"strconv"
	"strings"
)

type CreateInstutionResponse struct {
	InstitutionID string `json:"institution_id"`
	IsActive      bool   `json:"is_active"`
}

//! try to explain why first institution i take , what if some middleman just registers KU but doesnot register anything to

type CreateInstitutionRequest struct {
	InstitutionName string `json:"institution_name"`
	ToleAddress     string `json:"tole_address"`
	WardNumber      int    `json:"ward_number"`
	DistrictAddress string `json:"district_address"`
}
type CreateUserAccountRequest struct {
	InstitutionID         string `json:"institution_id"`
	Password              string `json:"password"`
	InstitutionRole       string `json:"institution_role"`
	SystemRole            string `json:"system_role"`
	InstitutionLogoBase64 string `json:"institution_logo_base64"`
	UserEmail             string `json:"user_email"`
}

type CreateUserAccountResponse struct {
	UserAccountID string `json:"user_acount_id"`
	CreatedAt     string `json:"created_at"`
}

func (c *CreateUserAccountRequest) ToEntity() entity.UserAccount {
	return entity.UserAccount{
		ID:              common.GenerateUUID(16),
		SystemRole:      enum.StringToRole(c.SystemRole),
		InstitutionRole: strings.ToUpper(c.InstitutionRole),
		Email:           c.UserEmail,
		Password:        c.Password,
	}

}

func (c *CreateInstitutionRequest) ToEntity() entity.Institution {
	return entity.Institution{
		InstitutionID:   common.GenerateUUID(16),
		InstitutionName: c.InstitutionName,
		ToleAddress:     c.ToleAddress,
		WardNumber:      strconv.Itoa(c.WardNumber),
		DistrictAddress: c.DistrictAddress,
		IsActive:        false,
	}
}
