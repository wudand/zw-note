package repository

import (
	"context"
	"database/sql"
	"fmt"

	"zw-note-backend/internal/model"

	"github.com/jmoiron/sqlx"
)

// RedemptionCodeRepository defines the data-access contract for redemption_codes.
type RedemptionCodeRepository interface {
	Create(ctx context.Context, rc *model.RedemptionCode) error
	// CreateBatch tx 可选；tx != nil 时参与外部事务，否则内部开启事务
	CreateBatch(ctx context.Context, tx *sqlx.Tx, codes []*model.RedemptionCode) error
	GetByID(ctx context.Context, id uint64) (*model.RedemptionCode, error)
	GetByCode(ctx context.Context, code string) (*model.RedemptionCode, error)
	List(ctx context.Context, offset, limit int, status *int8) ([]*model.RedemptionCode, int64, error)
	Update(ctx context.Context, rc *model.RedemptionCode) error
	Delete(ctx context.Context, id uint64) error
	// ClaimIfUnused 原子抢占：仅当 status=0 时更新，返回 true 表示抢成功
	ClaimIfUnused(ctx context.Context, code string, mpUserID uint64) (bool, error)
	// 用户最近 N 天内是否使用过任意兑换码
	HasUserUsedRecently(ctx context.Context, mpUserID uint64, withinDays int) (bool, error)
}

type redemptionCodeRepository struct {
	db *sqlx.DB
}

func NewRedemptionCodeRepository(db *sqlx.DB) RedemptionCodeRepository {
	return &redemptionCodeRepository{db: db}
}

func (r *redemptionCodeRepository) Create(ctx context.Context, rc *model.RedemptionCode) error {
	query := `INSERT INTO redemption_codes (code, status, mp_user_id, used_product_id, used_at)
		VALUES (:code, :status, :mp_user_id, :used_product_id, :used_at) RETURNING id`
	if err := namedReturningID(ctx, r.db, query, rc, &rc.ID); err != nil {
		return fmt.Errorf("create redemption_code: %w", err)
	}
	return nil
}

func (r *redemptionCodeRepository) CreateBatch(ctx context.Context, tx *sqlx.Tx, codes []*model.RedemptionCode) error {
	if len(codes) == 0 {
		return nil
	}
	ownTx := false
	if tx == nil {
		var err error
		tx, err = r.db.BeginTxx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin tx: %w", err)
		}
		defer tx.Rollback()
		ownTx = true
	}

	query := `INSERT INTO redemption_codes (code, status) VALUES ($1, $2) RETURNING id`
	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	for _, rc := range codes {
		rc.Status = model.RedemptionCodeStatusUnused
		if err := stmt.QueryRowxContext(ctx, rc.Code, rc.Status).Scan(&rc.ID); err != nil {
			return fmt.Errorf("insert: %w", err)
		}
	}

	if ownTx {
		return tx.Commit()
	}
	return nil
}

func (r *redemptionCodeRepository) GetByID(ctx context.Context, id uint64) (*model.RedemptionCode, error) {
	rc := &model.RedemptionCode{}
	err := r.db.GetContext(ctx, rc, `SELECT * FROM redemption_codes WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get redemption_code by id: %w", err)
	}
	return rc, nil
}

func (r *redemptionCodeRepository) GetByCode(ctx context.Context, code string) (*model.RedemptionCode, error) {
	rc := &model.RedemptionCode{}
	err := r.db.GetContext(ctx, rc, `SELECT * FROM redemption_codes WHERE code = $1`, code)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get redemption_code by code: %w", err)
	}
	return rc, nil
}

func (r *redemptionCodeRepository) List(ctx context.Context, offset, limit int, status *int8) ([]*model.RedemptionCode, int64, error) {
	countWhere := "TRUE"
	listWhere := "TRUE"
	args := []interface{}{}
	argN := 1
	if status != nil {
		countWhere += fmt.Sprintf(" AND status = $%d", argN)
		listWhere += fmt.Sprintf(" AND status = $%d", argN)
		args = append(args, *status)
		argN++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM redemption_codes WHERE %s", countWhere)
	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("count redemption_codes: %w", err)
	}

	listQuery := fmt.Sprintf(
		"SELECT * FROM redemption_codes WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		listWhere, argN, argN+1,
	)
	listArgs := append(args, limit, offset)
	var list []*model.RedemptionCode
	if err := r.db.SelectContext(ctx, &list, listQuery, listArgs...); err != nil {
		return nil, 0, fmt.Errorf("list redemption_codes: %w", err)
	}
	return list, total, nil
}

func (r *redemptionCodeRepository) Update(ctx context.Context, rc *model.RedemptionCode) error {
	query := `UPDATE redemption_codes SET status=:status, mp_user_id=:mp_user_id, used_product_id=:used_product_id, used_at=:used_at WHERE id=:id`
	if _, err := r.db.NamedExecContext(ctx, query, rc); err != nil {
		return fmt.Errorf("update redemption_code: %w", err)
	}
	return nil
}

// ClaimIfUnused 原子抢占：仅当 status=0 时更新，避免并发下多人同时抢到同一兑换码
func (r *redemptionCodeRepository) ClaimIfUnused(ctx context.Context, code string, mpUserID uint64) (bool, error) {
	result, err := r.db.ExecContext(ctx,
		`UPDATE redemption_codes SET status=$1, mp_user_id=$2, used_at=NOW() WHERE code=$3 AND status=$4`,
		model.RedemptionCodeStatusUsed, mpUserID, code, model.RedemptionCodeStatusUnused)
	if err != nil {
		return false, fmt.Errorf("claim redemption_code: %w", err)
	}
	rows, _ := result.RowsAffected()
	return rows == 1, nil
}

func (r *redemptionCodeRepository) Delete(ctx context.Context, id uint64) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM redemption_codes WHERE id = $1`, id); err != nil {
		return fmt.Errorf("delete redemption_code: %w", err)
	}
	return nil
}

func (r *redemptionCodeRepository) HasUserUsedRecently(ctx context.Context, mpUserID uint64, withinDays int) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM redemption_codes
		 WHERE mp_user_id = $1 AND status = 1 AND used_at > NOW() - make_interval(days => $2)`,
		mpUserID, withinDays)
	if err != nil {
		return false, fmt.Errorf("check user used recently: %w", err)
	}
	return count > 0, nil
}
