package repo

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/harry713j/vibe_writer/internal/model"
)

type RefreshTokenRepository struct {
	DB *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		DB: db,
	}
}

// create refresh token
func (r *RefreshTokenRepository) CreateRefreshToken(userId uuid.UUID) (*model.RefreshToken, error) {
	token := &model.RefreshToken{
		UserId:    userId,
		Token:     uuid.New(),
		CreatedAt: time.Now(),
		ExpireAt:  time.Now().Add(time.Hour * 24 * 7),
	}

	_, err := r.DB.Exec("INSERT INTO refresh_tokens(user_id, token, expire_at, created_at) VALUES($1, $2, $3, $4)",
		token.UserId, token.Token, token.CreatedAt, token.ExpireAt)

	if err != nil {
		return nil, err
	}

	return token, nil
}

// get refresh token
func (r *RefreshTokenRepository) GetRefreshToken(tokenValue uuid.UUID) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken

	err := r.DB.QueryRow("SELECT * FROM refresh_tokens WHERE token=$1", tokenValue).Scan(&refreshToken)

	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

// delete refresh token
func (r *RefreshTokenRepository) DeleteRefreshToken(userId uuid.UUID) error {
	_, err := r.DB.Exec("DELETE FROM refresh_tokens WHERE user_id=$1", userId)

	return err
}
