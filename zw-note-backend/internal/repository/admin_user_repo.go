package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"zw-note-backend/internal/model"

	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -destination=mock/admin_user_repo_mock.go -package=mock zw-note-backend/internal/repository AdminUserRepository

// AdminUserRepository defines the data-access contract for the admin_users table.
type AdminUserRepository interface {
	Create(ctx context.Context, user *model.AdminUser) error
	GetByID(ctx context.Context, id uint64) (*model.AdminUser, error)
	GetByUsername(ctx context.Context, username string) (*model.AdminUser, error)
	List(ctx context.Context, offset, limit int) ([]*model.AdminUser, int64, error)
	Update(ctx context.Context, user *model.AdminUser) error
	Delete(ctx context.Context, id uint64) error
}

type adminUserRepository struct {
	db *sqlx.DB
}

func NewAdminUserRepository(db *sqlx.DB) AdminUserRepository {
	return &adminUserRepository{db: db}
}

func (r *adminUserRepository) Create(ctx context.Context, user *model.AdminUser) error {
	query := `INSERT INTO admin_users (username, password, role, status)
		VALUES (:username, :password, :role, :status) RETURNING id`
	if err := namedReturningID(ctx, r.db, query, user, &user.ID); err != nil {
		return fmt.Errorf("create admin user: %w", err)
	}
	return nil
}

func (r *adminUserRepository) GetByID(ctx context.Context, id uint64) (*model.AdminUser, error) {
	user := &model.AdminUser{}
	err := r.db.GetContext(ctx, user, `SELECT * FROM admin_users WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get admin user by id: %w", err)
	}
	return user, nil
}

func (r *adminUserRepository) GetByUsername(ctx context.Context, username string) (*model.AdminUser, error) {
	user := &model.AdminUser{}
	err := r.db.GetContext(ctx, user, `SELECT * FROM admin_users WHERE username = $1`, username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get admin user by username: %w", err)
	}
	return user, nil
}

func (r *adminUserRepository) List(ctx context.Context, offset, limit int) ([]*model.AdminUser, int64, error) {
	var total int64
	if err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM admin_users`); err != nil {
		return nil, 0, fmt.Errorf("count admin users: %w", err)
	}

	var users []*model.AdminUser
	err := r.db.SelectContext(ctx, &users,
		`SELECT * FROM admin_users ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list admin users: %w", err)
	}
	return users, total, nil
}

func (r *adminUserRepository) Update(ctx context.Context, user *model.AdminUser) error {
	query := `UPDATE admin_users SET username=:username, role=:role, status=:status WHERE id=:id`
	if _, err := r.db.NamedExecContext(ctx, query, user); err != nil {
		return fmt.Errorf("update admin user: %w", err)
	}
	return nil
}

func (r *adminUserRepository) Delete(ctx context.Context, id uint64) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM admin_users WHERE id = $1`, id); err != nil {
		return fmt.Errorf("delete admin user: %w", err)
	}
	return nil
}
