package repository

import (
	"context"
	"database/sql"
	"fmt"

	"go-web-api/internal/model"

	"github.com/jmoiron/sqlx"
)

// AddressRepository defines the data-access contract for the addresses table.
type AddressRepository interface {
	Create(ctx context.Context, a *model.Address) error
	GetByID(ctx context.Context, id uint64) (*model.Address, error)
	GetByIDAndUser(ctx context.Context, id, mpUserID uint64) (*model.Address, error)
	ListByUser(ctx context.Context, mpUserID uint64) ([]*model.Address, error)
	CountByUser(ctx context.Context, mpUserID uint64) (int, error)
	// Update/ClearDefaultByUser/Delete 支持可选 tx，参与跨操作事务
	Update(ctx context.Context, tx *sqlx.Tx, a *model.Address) error
	ClearDefaultByUser(ctx context.Context, tx *sqlx.Tx, mpUserID uint64) error
	Delete(ctx context.Context, tx *sqlx.Tx, id uint64) error
}

type addressRepository struct {
	db *sqlx.DB
}

func NewAddressRepository(db *sqlx.DB) AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(ctx context.Context, a *model.Address) error {
	query := `INSERT INTO addresses (mp_user_id, receiver, phone, province, city, district, detail, tag, is_default)
		VALUES (:mp_user_id, :receiver, :phone, :province, :city, :district, :detail, :tag, :is_default)
		RETURNING id`
	if err := namedReturningID(ctx, r.db, query, a, &a.ID); err != nil {
		return fmt.Errorf("create address: %w", err)
	}
	return nil
}

func (r *addressRepository) GetByID(ctx context.Context, id uint64) (*model.Address, error) {
	a := &model.Address{}
	err := r.db.GetContext(ctx, a, `SELECT * FROM addresses WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get address by id: %w", err)
	}
	return a, nil
}

func (r *addressRepository) GetByIDAndUser(ctx context.Context, id, mpUserID uint64) (*model.Address, error) {
	a := &model.Address{}
	err := r.db.GetContext(ctx, a, `SELECT * FROM addresses WHERE id = $1 AND mp_user_id = $2`, id, mpUserID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get address by id and user: %w", err)
	}
	return a, nil
}

func (r *addressRepository) ListByUser(ctx context.Context, mpUserID uint64) ([]*model.Address, error) {
	var list []*model.Address
	err := r.db.SelectContext(ctx, &list,
		`SELECT * FROM addresses WHERE mp_user_id = $1 ORDER BY is_default DESC, created_at DESC`,
		mpUserID)
	if err != nil {
		return nil, fmt.Errorf("list addresses: %w", err)
	}
	return list, nil
}

func (r *addressRepository) CountByUser(ctx context.Context, mpUserID uint64) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM addresses WHERE mp_user_id = $1`, mpUserID)
	if err != nil {
		return 0, fmt.Errorf("count addresses: %w", err)
	}
	return count, nil
}

func (r *addressRepository) Update(ctx context.Context, tx *sqlx.Tx, a *model.Address) error {
	query := `UPDATE addresses SET receiver=:receiver, phone=:phone, province=:province, city=:city,
		district=:district, detail=:detail, tag=:tag, is_default=:is_default WHERE id=:id`
	if tx != nil {
		if _, err := tx.NamedExecContext(ctx, query, a); err != nil {
			return fmt.Errorf("update address: %w", err)
		}
	} else {
		if _, err := r.db.NamedExecContext(ctx, query, a); err != nil {
			return fmt.Errorf("update address: %w", err)
		}
	}
	return nil
}

func (r *addressRepository) ClearDefaultByUser(ctx context.Context, tx *sqlx.Tx, mpUserID uint64) error {
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, `UPDATE addresses SET is_default = 0 WHERE mp_user_id = $1`, mpUserID)
	} else {
		_, err = r.db.ExecContext(ctx, `UPDATE addresses SET is_default = 0 WHERE mp_user_id = $1`, mpUserID)
	}
	if err != nil {
		return fmt.Errorf("clear default addresses: %w", err)
	}
	return nil
}

func (r *addressRepository) Delete(ctx context.Context, tx *sqlx.Tx, id uint64) error {
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, `DELETE FROM addresses WHERE id = $1`, id)
	} else {
		_, err = r.db.ExecContext(ctx, `DELETE FROM addresses WHERE id = $1`, id)
	}
	if err != nil {
		return fmt.Errorf("delete address: %w", err)
	}
	return nil
}
