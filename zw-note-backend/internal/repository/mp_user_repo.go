package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go-web-api/internal/model"

	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -destination=mock/mp_user_repo_mock.go -package=mock go-web-api/internal/repository MPUserRepository

// MPUserRepository defines the data-access contract for the mp_users table.
type MPUserRepository interface {
	Create(ctx context.Context, user *model.MPUser) error
	GetByID(ctx context.Context, id uint64) (*model.MPUser, error)
	GetByOpenID(ctx context.Context, appID, openID string) (*model.MPUser, error)
	Update(ctx context.Context, user *model.MPUser) error
}

type mpUserRepository struct {
	db *sqlx.DB
}

func NewMPUserRepository(db *sqlx.DB) MPUserRepository {
	return &mpUserRepository{db: db}
}

func (r *mpUserRepository) Create(ctx context.Context, user *model.MPUser) error {
	query := `INSERT INTO mp_users (app_id, openid, unionid, nickname, avatar, status)
	          VALUES (:app_id, :openid, :unionid, :nickname, :avatar, :status) RETURNING id`
	if err := namedReturningID(ctx, r.db, query, user, &user.ID); err != nil {
		return fmt.Errorf("create mp user: %w", err)
	}
	return nil
}

func (r *mpUserRepository) GetByID(ctx context.Context, id uint64) (*model.MPUser, error) {
	user := &model.MPUser{}
	err := r.db.GetContext(ctx, user, `SELECT * FROM mp_users WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get mp user by id: %w", err)
	}
	return user, nil
}

// GetByOpenID looks up a mini-program user by the composite (app_id, openid) key.
func (r *mpUserRepository) GetByOpenID(ctx context.Context, appID, openID string) (*model.MPUser, error) {
	user := &model.MPUser{}
	err := r.db.GetContext(ctx, user,
		`SELECT * FROM mp_users WHERE app_id = $1 AND openid = $2`, appID, openID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get mp user by openid: %w", err)
	}
	return user, nil
}

func (r *mpUserRepository) Update(ctx context.Context, user *model.MPUser) error {
	query := `UPDATE mp_users SET nickname=:nickname, avatar=:avatar, status=:status WHERE id=:id`
	if _, err := r.db.NamedExecContext(ctx, query, user); err != nil {
		return fmt.Errorf("update mp user: %w", err)
	}
	return nil
}
