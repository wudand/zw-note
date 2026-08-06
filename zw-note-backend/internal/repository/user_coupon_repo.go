package repository

import (
	"context"
	"database/sql"
	"fmt"

	"zw-note-backend/internal/model"

	"github.com/jmoiron/sqlx"
)

// UserCouponRepository defines the data-access contract for the user_coupons table.
type UserCouponRepository interface {
	Create(ctx context.Context, uc *model.UserCoupon) error
	GetByID(ctx context.Context, id uint64) (*model.UserCoupon, error)
	GetByIDAndUser(ctx context.Context, id, mpUserID uint64) (*model.UserCoupon, error)
	ListByUser(ctx context.Context, mpUserID uint64, status *string) ([]*model.UserCoupon, error)
	CountUnusedByUserAndCoupon(ctx context.Context, mpUserID, couponID uint64) (int, error)
	HasUsedNewUserCoupon(ctx context.Context, mpUserID uint64) (bool, error)
	Update(ctx context.Context, uc *model.UserCoupon) error
}

type userCouponRepository struct {
	db *sqlx.DB
}

func NewUserCouponRepository(db *sqlx.DB) UserCouponRepository {
	return &userCouponRepository{db: db}
}

func (r *userCouponRepository) Create(ctx context.Context, uc *model.UserCoupon) error {
	query := `INSERT INTO user_coupons (mp_user_id, coupon_id, coupon_name, coupon_type, min_amount, discount_value, stackable, status, claimed_at, expiry_at, used_at)
		VALUES (:mp_user_id, :coupon_id, :coupon_name, :coupon_type, :min_amount, :discount_value, :stackable, :status, :claimed_at, :expiry_at, :used_at)
		RETURNING id`
	if err := namedReturningID(ctx, r.db, query, uc, &uc.ID); err != nil {
		return fmt.Errorf("create user_coupon: %w", err)
	}
	return nil
}

func (r *userCouponRepository) GetByID(ctx context.Context, id uint64) (*model.UserCoupon, error) {
	uc := &model.UserCoupon{}
	err := r.db.GetContext(ctx, uc, `SELECT * FROM user_coupons WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user_coupon by id: %w", err)
	}
	return uc, nil
}

func (r *userCouponRepository) GetByIDAndUser(ctx context.Context, id, mpUserID uint64) (*model.UserCoupon, error) {
	uc := &model.UserCoupon{}
	err := r.db.GetContext(ctx, uc, `SELECT * FROM user_coupons WHERE id = $1 AND mp_user_id = $2`, id, mpUserID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user_coupon by id and user: %w", err)
	}
	return uc, nil
}

func (r *userCouponRepository) ListByUser(ctx context.Context, mpUserID uint64, status *string) ([]*model.UserCoupon, error) {
	query := `SELECT * FROM user_coupons WHERE mp_user_id = $1`
	args := []interface{}{mpUserID}
	if status != nil {
		query += ` AND status = $2`
		args = append(args, *status)
	}
	query += ` ORDER BY expiry_at ASC`

	var list []*model.UserCoupon
	if err := r.db.SelectContext(ctx, &list, query, args...); err != nil {
		return nil, fmt.Errorf("list user_coupons: %w", err)
	}
	return list, nil
}

func (r *userCouponRepository) CountUnusedByUserAndCoupon(ctx context.Context, mpUserID, couponID uint64) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM user_coupons WHERE mp_user_id = $1 AND coupon_id = $2 AND status = 'unused'`,
		mpUserID, couponID)
	if err != nil {
		return 0, fmt.Errorf("count unused user_coupons: %w", err)
	}
	return count, nil
}

func (r *userCouponRepository) HasUsedNewUserCoupon(ctx context.Context, mpUserID uint64) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM user_coupons WHERE mp_user_id = $1 AND coupon_type = 'new_user' AND status = 'used'`,
		mpUserID)
	if err != nil {
		return false, fmt.Errorf("check used new_user coupon: %w", err)
	}
	return count > 0, nil
}

func (r *userCouponRepository) Update(ctx context.Context, uc *model.UserCoupon) error {
	query := `UPDATE user_coupons SET status=:status, used_at=:used_at WHERE id=:id`
	if _, err := r.db.NamedExecContext(ctx, query, uc); err != nil {
		return fmt.Errorf("update user_coupon: %w", err)
	}
	return nil
}
