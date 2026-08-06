package repository

import (
	"context"
	"database/sql"
	"fmt"

	"zw-note-backend/internal/model"

	"github.com/jmoiron/sqlx"
)

// CouponRepository defines the data-access contract for the coupons table.
type CouponRepository interface {
	Create(ctx context.Context, c *model.Coupon) error
	GetByID(ctx context.Context, id uint64) (*model.Coupon, error)
	List(ctx context.Context, offset, limit int, couponType *string, status *int8) ([]*model.Coupon, int64, error)
	Update(ctx context.Context, c *model.Coupon) error
	Delete(ctx context.Context, id uint64) error
}

type couponRepository struct {
	db *sqlx.DB
}

func NewCouponRepository(db *sqlx.DB) CouponRepository {
	return &couponRepository{db: db}
}

func (r *couponRepository) Create(ctx context.Context, c *model.Coupon) error {
	query := `INSERT INTO coupons (name, type, min_amount, discount_value, valid_days, stackable, status)
		VALUES (:name, :type, :min_amount, :discount_value, :valid_days, :stackable, :status)
		RETURNING id`
	if err := namedReturningID(ctx, r.db, query, c, &c.ID); err != nil {
		return fmt.Errorf("create coupon: %w", err)
	}
	return nil
}

func (r *couponRepository) GetByID(ctx context.Context, id uint64) (*model.Coupon, error) {
	c := &model.Coupon{}
	err := r.db.GetContext(ctx, c, `SELECT * FROM coupons WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get coupon by id: %w", err)
	}
	return c, nil
}

func (r *couponRepository) List(ctx context.Context, offset, limit int, couponType *string, status *int8) ([]*model.Coupon, int64, error) {
	countWhere := "TRUE"
	listWhere := "TRUE"
	args := []interface{}{}
	argN := 1
	if couponType != nil {
		countWhere += fmt.Sprintf(" AND type = $%d", argN)
		listWhere += fmt.Sprintf(" AND type = $%d", argN)
		args = append(args, *couponType)
		argN++
	}
	if status != nil {
		countWhere += fmt.Sprintf(" AND status = $%d", argN)
		listWhere += fmt.Sprintf(" AND status = $%d", argN)
		args = append(args, *status)
		argN++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM coupons WHERE %s", countWhere)
	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("count coupons: %w", err)
	}

	listQuery := fmt.Sprintf(
		"SELECT * FROM coupons WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		listWhere, argN, argN+1,
	)
	listArgs := append(args, limit, offset)
	var list []*model.Coupon
	if err := r.db.SelectContext(ctx, &list, listQuery, listArgs...); err != nil {
		return nil, 0, fmt.Errorf("list coupons: %w", err)
	}
	return list, total, nil
}

func (r *couponRepository) Update(ctx context.Context, c *model.Coupon) error {
	query := `UPDATE coupons SET name=:name, min_amount=:min_amount, discount_value=:discount_value, valid_days=:valid_days, status=:status WHERE id=:id`
	if _, err := r.db.NamedExecContext(ctx, query, c); err != nil {
		return fmt.Errorf("update coupon: %w", err)
	}
	return nil
}

func (r *couponRepository) Delete(ctx context.Context, id uint64) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM coupons WHERE id = $1`, id); err != nil {
		return fmt.Errorf("delete coupon: %w", err)
	}
	return nil
}
