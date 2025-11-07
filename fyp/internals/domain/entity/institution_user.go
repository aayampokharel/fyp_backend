package entity

type InstitutionUser struct {
	InstitutionID         string `json:"institution_id"`
	UserID                string `json:"user_id"`
	InstitutionLogoBase64 string `json:"institution_logo_base64"`
}
