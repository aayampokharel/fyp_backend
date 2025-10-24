package entity

import (
	"project/package/enum"
	"time"
)

type UserAccount struct {
	ID              string     `json:"id"`
	SystemRole      enum.ROLE  `json:"system_role"`
	InstitutionRole string     `json:"institution_role"`
	CreatedAt       time.Time  `json:"created_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
	Email           string     `json:"email"`
	Password        string     `json:"password"`
}
