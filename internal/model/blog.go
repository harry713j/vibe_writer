package model

import (
	"time"

	"github.com/google/uuid"
)

type Blog struct {
	Id         int64      `json:"id"`
	UserId     uuid.UUID  `json:"user_id"`
	Title      string     `json:"title"`
	Slug       string     `json:"slug"`
	Content    string     `json:"content"`
	Visibility bool       `json:"visibility"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type BlogPhoto struct {
	Id        int64     `json:"id"`
	BlogId    int64     `json:"blog_id"`
	PhotoUrl  string    `json:"photo_url"`
	CreatedAt time.Time `json:"created_at"`
}

type BlogRes struct {
	Blog
	PhotoUrls []string `json:"photo_urls"`
}
