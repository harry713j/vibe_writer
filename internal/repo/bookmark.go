package repo

import (
	"database/sql"

	"github.com/google/uuid"
)

type BookmarkRepository struct {
	DB *sql.DB
}

func NewBookmarkRepository(db *sql.DB) *BookmarkRepository {
	return &BookmarkRepository{DB: db}
}

func (b *BookmarkRepository) Upsert(userId uuid.UUID, blogId int64) error {
	query := `
		INSERT INTO bookmarks(user_id, blog_id) VALUES($1, $2)
		ON CONFLICT(user_id, blog_id) DO NOTHNG
	`
	if _, err := b.DB.Exec(query, userId, blogId); err != nil {
		return err
	}

	return nil
}

func (b *BookmarkRepository) Delete(userId uuid.UUID, blogId int64) error {
	query := `
		DELETE FROM bookmarks WHERE user_id = $1 AND blog_id = $2
	`

	if _, err := b.DB.Exec(query, userId, blogId); err != nil {
		return err
	}
	return nil
}
