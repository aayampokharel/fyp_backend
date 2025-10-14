package entity

type Institution struct {
	InstitutionID   string `json:"institution_id"`
	InstitutionName string `json:"institution_name"`
	ToleAddress     string `json:"tole_address"`
	DistrictAddress string `json:"district_address"`
	IsActive        bool   `json:"is_active"`
}
