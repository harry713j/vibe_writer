package repo

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/harry713j/vibe_writer/internal/model"
)

type CommentRepository struct {
	DB *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{
		DB: db,
	}
}

// create a comment
func (c *CommentRepository) CreateComment(userId uuid.UUID, blogId int64, parentId int64, content string) (int64, error) {
	var query string
	var commentId int64

	if parentId != 0 {
		query = `INSERT INTO comments(user_id, blog_id, parent_id, content)
			VALUES($1, $2, $3, $4) RETURNING id`

		err := c.DB.QueryRow(query, userId, blogId, parentId, content).Scan(&commentId)

		if err != nil {
			return 0, err
		}

		return commentId, nil
	} else {
		query = `INSERT INTO comments(user_id, blog_id, content)
			VALUES($1, $2, $3) RETURNING id`

		err := c.DB.QueryRow(query, userId, blogId, content).Scan(&commentId)

		if err != nil {
			return 0, err
		}

		return commentId, nil
	}
}

// Get a comment
func (c *CommentRepository) GetCommentById(userId uuid.UUID, id int64) (*model.CommentWithStat, error) {
	var comment model.CommentWithStat

	query := `
		SELECT
			c.id,
			c.user_id,
			c.parent_id,
			c.content,
			c.created_at,
			c.updated_at,
			COUNT(l.id) FILTER (WHERE l.like_type = 'like') AS likes_count,
			COUNT(l.id) FILTER (WHERE l.like_type = 'dislike') AS dislikes_count
		FROM comments c
		LEFT JOIN likes l ON c.id = l.comment_id
		WHERE c.user_id = $1 AND c.id = $2
		GROUP BY c.id
	`

	err := c.DB.QueryRow(query, userId, id).Scan(
		&comment.Id, &comment.UserId, &comment.ParentId, &comment.Content,
		&comment.CreatedAt, &comment.UpdatedAt, &comment.LikeCount, &comment.DislikeCount,
	)

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// Get comments of a blog
func (c *CommentRepository) GetCommentsByBlogId(blogId int64) ([]model.CommentWithStat, error) {
	var comments []model.CommentWithStat

	query := `
		SELECT
			c.id,
			c.user_id,
			c.parent_id,
			c.content,
			c.created_at,
			c.updated_at,
			COUNT(l.id) FILTER (WHERE l.like_type = 'like') AS likes_count,
			COUNT(l.id) FILTER (WHERE l.like_type = 'dislike') AS dislikes_count
		FROM comments c
		LEFT JOIN likes l ON c.id = l.comment_id
		WHERE c.blog_id = $1
		GROUP BY c.id
		ORDER BY c.created_at DESC
	`

	rows, err := c.DB.Query(query, blogId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		comment := model.CommentWithStat{}

		err := rows.Scan(
			&comment.Id, &comment.UserId, &comment.ParentId, &comment.Content,
			&comment.CreatedAt, &comment.UpdatedAt, &comment.LikeCount, &comment.DislikeCount,
		)

		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

// delete a comment
func (c *CommentRepository) DeleteCommentById(userId uuid.UUID, id int64) error {
	if _, err := c.DB.Exec("DELETE FROM comments WHERE id=$1 AND user_id=$2", id, userId); err != nil {
		return err
	}

	return nil
}
