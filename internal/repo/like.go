package repo

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/harry713j/vibe_writer/internal/model"
)

type LikeRepository struct {
	DB *sql.DB
}

func NewLikeRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{DB: db}
}

// insert and update
func (r *LikeRepository) UpsertCommentLike(userId uuid.UUID, commentId int64, liketype model.LikeType) (*model.Like, error) {
	var like model.Like

	err := r.DB.QueryRow(`INSERT INTO likes(user_id, comment_id, liketype) VALUES($1, $2, $3)
		ON CONFLICT(user_id, comment_id) 
		DO UPDATE SET
		like_type = EXCLUDED.like_type,
		updated_at = CURRENT_TIMESTAMP
		  RETURNNING *`,
		userId, commentId, liketype).Scan(
		&like.Id, &like.UserId, &like.BlogId, &like.CommentId, &like.LikeType,
		&like.CreatedAt, &like.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &like, nil
}

func (r *LikeRepository) UpsertBlogLike(userId uuid.UUID, blogId int64, liketype model.LikeType) (*model.Like, error) {
	var like model.Like

	err := r.DB.QueryRow(`INSERT INTO likes(user_id, blog_id, liketype) VALUES($1, $2, $3) 
			ON CONFLICT(user_id, blog_id)
			DO UPDATE SET
			like_type = EXCLUDED.like_type,
			updated_at = CURRENT_TIMESTAMP
			RETURNNING *`,
		userId, blogId, liketype).Scan(
		&like.Id, &like.UserId, &like.BlogId, &like.CommentId, &like.LikeType,
		&like.CreatedAt, &like.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &like, nil
}

// delete
func (r *LikeRepository) DeleteCommentLike(userId uuid.UUID, commentId int64) error {
	if _, err := r.DB.Exec("DELETE FROM likes WHERE user_id=$1 AND comment_id=$2", userId, commentId); err != nil {
		return err
	}

	return nil
}

func (r *LikeRepository) DeleteBlogLike(userId uuid.UUID, blogId int64) error {
	if _, err := r.DB.Exec("DELETE FROM likes WHERE user_id=$1 AND blog_id=$2", userId, blogId); err != nil {
		return err
	}

	return nil
}
