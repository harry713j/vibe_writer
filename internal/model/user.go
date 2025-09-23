package model

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	UserId    uuid.UUID `json:"user_id"`
	FullName  string    `json:"full_name"`
	Bio       string    `json:"bio"`
	AvatarUrl string    `json:"avatar_url"`
}

type UserDetails struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FullName  string    `json:"full_name"`
	Bio       string    `json:"bio"`
	AvatarUrl string    `json:"avatar_url"`
}
