package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// RedemptionCodeProductRepository defines the data-access contract for redemption_code_products.
type RedemptionCodeProductRepository interface {
	CreateBatch(ctx context.Context, tx *sqlx.Tx, redemptionCodeID uint64, productIDs []uint64) error
	GetProductIDsByCodeID(ctx context.Context, redemptionCodeID uint64) ([]uint64, error)
	DeleteByCodeID(ctx context.Context, tx *sqlx.Tx, redemptionCodeID uint64) error
	HasProduct(ctx context.Context, redemptionCodeID, productID uint64) (bool, error)
}

type redemptionCodeProductRepository struct {
	db *sqlx.DB
}

func NewRedemptionCodeProductRepository(db *sqlx.DB) RedemptionCodeProductRepository {
	return &redemptionCodeProductRepository{db: db}
}

func (r *redemptionCodeProductRepository) CreateBatch(ctx context.Context, tx *sqlx.Tx, redemptionCodeID uint64, productIDs []uint64) error {
	if len(productIDs) == 0 {
		return nil
	}
	query := `INSERT INTO redemption_code_products (redemption_code_id, product_id) VALUES ($1, $2)`
	for _, pid := range productIDs {
		var err error
		if tx != nil {
			_, err = tx.ExecContext(ctx, query, redemptionCodeID, pid)
		} else {
			_, err = r.db.ExecContext(ctx, query, redemptionCodeID, pid)
		}
		if err != nil {
			return fmt.Errorf("create redemption_code_product: %w", err)
		}
	}
	return nil
}

func (r *redemptionCodeProductRepository) GetProductIDsByCodeID(ctx context.Context, redemptionCodeID uint64) ([]uint64, error) {
	var ids []uint64
	err := r.db.SelectContext(ctx, &ids,
		`SELECT product_id FROM redemption_code_products WHERE redemption_code_id = $1 ORDER BY id`,
		redemptionCodeID)
	if err != nil {
		return nil, fmt.Errorf("get product ids by code id: %w", err)
	}
	return ids, nil
}

func (r *redemptionCodeProductRepository) DeleteByCodeID(ctx context.Context, tx *sqlx.Tx, redemptionCodeID uint64) error {
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, `DELETE FROM redemption_code_products WHERE redemption_code_id = $1`, redemptionCodeID)
	} else {
		_, err = r.db.ExecContext(ctx, `DELETE FROM redemption_code_products WHERE redemption_code_id = $1`, redemptionCodeID)
	}
	if err != nil {
		return fmt.Errorf("delete redemption_code_products: %w", err)
	}
	return nil
}

func (r *redemptionCodeProductRepository) HasProduct(ctx context.Context, redemptionCodeID, productID uint64) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM redemption_code_products WHERE redemption_code_id = $1 AND product_id = $2`,
		redemptionCodeID, productID)
	if err != nil {
		return false, fmt.Errorf("check has product: %w", err)
	}
	return count > 0, nil
}
