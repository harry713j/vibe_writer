package model

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	Id        int64      `json:"id"`
	UserId    uuid.UUID  `json:"user_id"`
	BlogId    int64      `json:"blog_id"`
	ParentId  int64      `json:"parent_id"`
	Content   string     `json:"content"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
