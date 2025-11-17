package entity

import "time"

type Institution struct {
	InstitutionID     string     `json:"institution_id"`
	InstitutionName   string     `json:"institution_name"`
	ToleAddress       string     `json:"tole_address"`
	WardNumber        string     `json:"ward_number"`
	DistrictAddress   string     `json:"district_address"`
	IsActive          *bool      `json:"is_active"`
	IsSignupCompleted bool       `json:"is_signup_completed"`
	CreatedAt         time.Time  `json:"created_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
}
