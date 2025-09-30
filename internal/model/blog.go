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

type BlogSummary struct {
	Blog
	Thumbnail    string `json:"blog_thumbnail"`
	LikesCount   int    `json:"likes_count"`
	DislikeCount int    `json:"dislikes_count"`
	CommentCount int    `json:"comments_count"`
}

type BlogWithStat struct {
	Blog
	PhotoUrls    []string `json:"photo_urls"`
	LikeCount    int      `json:"likes_count"`
	DislikeCount int      `json:"dislikes_count"`
}

type BlogResponse struct {
	BlogWithStat
	Comments     []CommentWithStat `json:"comments"`
	AuthorName   string            `json:"author_name"`
	AuthorBio    string            `json:"author_bio"`
	AuthorAvatar string            `json:"author_avatar"`
}
