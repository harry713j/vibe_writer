package model

import (
	"time"

	"github.com/google/uuid"
)

// enum type
type LikeType string

const (
	LIKE    LikeType = "like"
	DISLIKE LikeType = "dislike"
)

type Like struct {
	Id        int64      `json:"id"`
	UserId    uuid.UUID  `json:"user_id"`
	BlogId    int64      `json:"blog_id"`
	CommentId int64      `json:"comment_id"`
	LikeType  LikeType   `json:"like_type"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
