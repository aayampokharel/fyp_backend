package models

import (
	"project/internals/domain/entity"
	"time"
)

type UserAccount struct {
	ID              string     `json:"id"`
	SystemRole      string     `json:"system_role"`
	InstitutionRole string     `json:"institution_role"`
	CreatedAt       time.Time  `json:"created_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
	Email           string     `json:"email"`
	Password        string     `json:"password"`
}

func UserAccountFromEntity(userAccount entity.UserAccount) UserAccount {
	return UserAccount{
		ID:              userAccount.ID,
		SystemRole:      userAccount.SystemRole.String(),
		Email:           userAccount.Email,
		Password:        userAccount.Password,
		InstitutionRole: userAccount.InstitutionRole,
	}
}
