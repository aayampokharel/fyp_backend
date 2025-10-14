package entity

import "time"

type UserAccount struct {
	ID        string     `json:"id"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
}
