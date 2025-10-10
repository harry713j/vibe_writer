package model

import (
	"time"

	"github.com/google/uuid"
)

type Bookmark struct {
	UserId    uuid.UUID  `json:"user_id"`
	BlogId    int64      `json:"blog_id"`
	CreatedAt *time.Time `json:"created_at"`
}
