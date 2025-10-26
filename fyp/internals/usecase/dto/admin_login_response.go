package dto

import (
	"project/internals/domain/entity"
	"time"
)

type AdminLoginResponse struct {
	UserID               string               `json:"user_id"`
	GeneratedUniqueToken string               `json:"generated_unique_token"`
	CreatedTime          time.Time            `json:"created_time"`
	InstitutionList      []entity.Institution `json:"institution_list"`
}
