package repository

import (
	"context"
	"database/sql"
	"fmt"

	"go-web-api/internal/model"

	"github.com/jmoiron/sqlx"
)

// CategoryRepository defines the data-access contract for the product_categories table.
type CategoryRepository interface {
	Create(ctx context.Context, cat *model.Category) error
	GetByID(ctx context.Context, id uint64) (*model.Category, error)
	List(ctx context.Context, offset, limit int, includeDisabled bool) ([]*model.Category, int64, error)
	Update(ctx context.Context, cat *model.Category) error
	Delete(ctx context.Context, id uint64) error
}

type categoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, cat *model.Category) error {
	query := `INSERT INTO product_categories (name, sort_order, status)
		VALUES (:name, :sort_order, :status) RETURNING id`
	if err := namedReturningID(ctx, r.db, query, cat, &cat.ID); err != nil {
		return fmt.Errorf("create category: %w", err)
	}
	return nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id uint64) (*model.Category, error) {
	cat := &model.Category{}
	err := r.db.GetContext(ctx, cat, `SELECT * FROM product_categories WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get category by id: %w", err)
	}
	return cat, nil
}

func (r *categoryRepository) List(ctx context.Context, offset, limit int, includeDisabled bool) ([]*model.Category, int64, error) {
	countQuery := `SELECT COUNT(*) FROM product_categories`
	if !includeDisabled {
		countQuery = `SELECT COUNT(*) FROM product_categories WHERE status = 1`
	}
	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, fmt.Errorf("count categories: %w", err)
	}

	var list []*model.Category
	if includeDisabled {
		err := r.db.SelectContext(ctx, &list,
			`SELECT * FROM product_categories ORDER BY sort_order ASC, id ASC LIMIT $1 OFFSET $2`, limit, offset)
		if err != nil {
			return nil, 0, fmt.Errorf("list categories: %w", err)
		}
	} else {
		err := r.db.SelectContext(ctx, &list,
			`SELECT * FROM product_categories WHERE status = 1 ORDER BY sort_order ASC, id ASC LIMIT $1 OFFSET $2`, limit, offset)
		if err != nil {
			return nil, 0, fmt.Errorf("list categories: %w", err)
		}
	}
	return list, total, nil
}

func (r *categoryRepository) Update(ctx context.Context, cat *model.Category) error {
	query := `UPDATE product_categories SET name=:name, sort_order=:sort_order, status=:status WHERE id=:id`
	if _, err := r.db.NamedExecContext(ctx, query, cat); err != nil {
		return fmt.Errorf("update category: %w", err)
	}
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id uint64) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM product_categories WHERE id = $1`, id); err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	return nil
}
