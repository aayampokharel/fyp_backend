package admin

import "project/internals/domain/entity"

type AdminLoginRequest struct {
	AdminEmail string `json:"admin_email"`
	Password   string `json:"password"`
}

type AdminLoginResponse struct {
	SSEToken        string               `json:"sse_token"`
	UserID          string               `json:"user_id"`
	CreatedAt       string               `json:"created_at"`
	InstitutionList []entity.Institution `json:"institution_list"`

	////many more nested maps
}
