package entity

type InstitutionUser struct {
	InstitutionID            string `json:"institution_id"`
	UserID                   string `json:"user_id"`
	PublicKey                string `json:"public_key"`
	PrincipalName            string `json:"principal_name"`
	PrincipalSignatureBase64 string `json:"principal_signature_base64"`
	InstitutionLogoBase64    string `json:"institution_logo_base64"`
}
