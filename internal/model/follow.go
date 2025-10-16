package model

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	FollowerId  uuid.UUID  `json:"follower_id"`
	FollowingId uuid.UUID  `json:"following_id"`
	CreatedAt   *time.Time `json:"created_at"`
}

type FollowResponse struct {
	UserId   uuid.UUID `json:"user_id"`
	FullName string    `json:"full_name"`
	Bio      string    `json:"bio"`
	Avatar   string    `json:"avatar_url"`
}
