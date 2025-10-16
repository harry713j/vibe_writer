package repo

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/harry713j/vibe_writer/internal/model"
)

type FollowRepository struct {
	DB *sql.DB
}

func NewFollowRepository(db *sql.DB) *FollowRepository {
	return &FollowRepository{DB: db}
}

func (f *FollowRepository) Create(followerId uuid.UUID, followingId uuid.UUID) error {
	if _, err := f.DB.Exec("INSERT INTO follows(follower_id, following_id) VALUES($1, $2)", followerId, followingId); err != nil {
		return err
	}

	return nil
}

func (f *FollowRepository) Delete(followerId uuid.UUID, followingId uuid.UUID) error {
	if _, err := f.DB.Exec("DELETE FROM follows WHERE follower_id = $1 AND following_id = $2", followerId, followingId); err != nil {
		return err
	}

	return nil
}

func (f *FollowRepository) GetAllFollower(followingId uuid.UUID, page, limit int) (*model.PaginatedResponse[model.FollowResponse], error) {
	if page < 1 {
		page = 1
	}

	if limit <= 0 {
		limit = 20
	}

	offset := (page - 1) * limit

	query := `
		SELECT 
			up.full_name,
			up.bio,
			up.avatar_url,
			up.user_id
		FROM follows f
		JOIN user_profiles up ON f.follower_id = up.user_id
		WHERE f.following_id = $1
		LIMIT $2 OFFSET $3
	`

	var followers []model.FollowResponse

	rows, err := f.DB.Query(query, followingId, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var follower model.FollowResponse
		err := rows.Scan(
			&follower.FullName, &follower.Bio, &follower.Avatar, &follower.UserId,
		)

		if err != nil {
			return nil, err
		}

		followers = append(followers, follower)
	}

	var total int
	err = f.DB.QueryRow("SELECT COUNT(*) FROM follows WHERE following_id = $1", followingId).Scan(&total)
	if err != nil {
		return nil, err
	}

	totalPages := (total + limit - 1) / limit

	return &model.PaginatedResponse[model.FollowResponse]{
		Data: followers,
		Meta: model.PageMeta{
			Total: total,
			Pages: totalPages,
			Page:  page,
			Limit: limit,
		},
	}, nil
}

func (f *FollowRepository) GetAllFollowing(followerId uuid.UUID, page, limit int) (*model.PaginatedResponse[model.FollowResponse], error) {
	if page < 1 {
		page = 1
	}

	if limit <= 0 {
		limit = 20
	}

	offset := (page - 1) * limit

	query := `
		SELECT 
			up.full_name,
			up.bio,
			up.avatar_url,
			up.user_id
		FROM follows f
		JOIN user_profiles up ON f.following_id = up.user_id
		WHERE f.follower_id = $1
		LIMIT $2 OFFSET $3
	`

	var followings []model.FollowResponse

	rows, err := f.DB.Query(query, followerId, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var following model.FollowResponse
		err := rows.Scan(
			&following.FullName, &following.Bio, &following.Avatar, &following.UserId,
		)

		if err != nil {
			return nil, err
		}

		followings = append(followings, following)
	}

	var total int
	err = f.DB.QueryRow("SELECT COUNT(*) FROM follows WHERE follower_id = $1", followerId).Scan(&total)
	if err != nil {
		return nil, err
	}

	totalPages := (total + limit - 1) / limit

	return &model.PaginatedResponse[model.FollowResponse]{
		Data: followings,
		Meta: model.PageMeta{
			Total: total,
			Pages: totalPages,
			Page:  page,
			Limit: limit,
		},
	}, nil
}
