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
func (c *CommentRepository) CreateComment(userId uuid.UUID, blogId int64, parentId int64, content string) (*model.Comment, error) {
	var query string
	var comment model.Comment

	if parentId != 0 {
		query = `INSERT INTO comments(user_id, blog_id, parent_id, content) 
			VALUES($1, $2, $3, $4) RETURNING id`

		err := c.DB.QueryRow(query, userId, blogId, parentId, content).Scan(&comment)

		if err != nil {
			return nil, err
		}

		return &comment, nil
	} else {
		query = `INSERT INTO comments(user_id, blog_id, content) 
			VALUES($1, $2, $3) RETURNNING id`

		err := c.DB.QueryRow(query, userId, blogId, content).Scan(&comment)

		if err != nil {
			return nil, err
		}

		return &comment, nil
	}
}

// Get a comment
func (c *CommentRepository) GetCommentById(userId uuid.UUID, id int64) (*model.Comment, error) {
	var comment model.Comment

	err := c.DB.QueryRow("SELECT * FROM comments WHERE user_id=$1 AND blog_id=$2", userId, id).Scan(&comment)

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (c *CommentRepository) GetComment(id int64) (*model.Comment, error) {
	var comment model.Comment

	err := c.DB.QueryRow("SELECT * FROM comments WHERE AND blog_id=$2", id).Scan(&comment)

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// Get comments of a blog
func (c *CommentRepository) GetCommentsByBlogId(blogId int64) ([]*model.Comment, error) {
	var comments []*model.Comment

	rows, err := c.DB.Query("SELECT * FROM comments WHERE blog_id=$1", blogId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		comment := &model.Comment{}

		err := rows.Scan(comment)

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
