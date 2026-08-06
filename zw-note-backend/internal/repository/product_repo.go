package repository

import (
	"context"
	"database/sql"
	"fmt"

	"zw-note-backend/internal/model"

	"github.com/jmoiron/sqlx"
)

// ProductRepository defines the data-access contract for the products table.
type ProductRepository interface {
	Create(ctx context.Context, p *model.Product) error
	GetByID(ctx context.Context, id uint64) (*model.Product, error)
	List(ctx context.Context, offset, limit int, categoryID *uint64, includeDisabled bool) ([]*model.Product, int64, error)
	Update(ctx context.Context, p *model.Product) error
	Delete(ctx context.Context, id uint64) error
}

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, p *model.Product) error {
	query := `INSERT INTO products (category_id, name, ingredients, grade, storage, specification, cover_image, carousel_images, detail, detail_images, status)
		VALUES (:category_id, :name, :ingredients, :grade, :storage, :specification, :cover_image, :carousel_images, :detail, :detail_images, :status)
		RETURNING id`
	if err := namedReturningID(ctx, r.db, query, p, &p.ID); err != nil {
		return fmt.Errorf("create product: %w", err)
	}
	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id uint64) (*model.Product, error) {
	p := &model.Product{}
	err := r.db.GetContext(ctx, p, `SELECT * FROM products WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get product by id: %w", err)
	}
	return p, nil
}

func (r *productRepository) List(ctx context.Context, offset, limit int, categoryID *uint64, includeDisabled bool) ([]*model.Product, int64, error) {
	countWhere := "TRUE"
	listWhere := "TRUE"
	args := []interface{}{}
	argN := 1
	if !includeDisabled {
		countWhere += " AND status = 1"
		listWhere += " AND status = 1"
	}
	if categoryID != nil {
		countWhere += fmt.Sprintf(" AND category_id = $%d", argN)
		listWhere += fmt.Sprintf(" AND category_id = $%d", argN)
		args = append(args, *categoryID)
		argN++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products WHERE %s", countWhere)
	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}

	listQuery := fmt.Sprintf(
		"SELECT * FROM products WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		listWhere, argN, argN+1,
	)
	listArgs := append(args, limit, offset)
	var list []*model.Product
	if err := r.db.SelectContext(ctx, &list, listQuery, listArgs...); err != nil {
		return nil, 0, fmt.Errorf("list products: %w", err)
	}
	return list, total, nil
}

func (r *productRepository) Update(ctx context.Context, p *model.Product) error {
	query := `UPDATE products SET category_id=:category_id, name=:name, ingredients=:ingredients, grade=:grade, storage=:storage,
		specification=:specification, cover_image=:cover_image, carousel_images=:carousel_images, detail=:detail, detail_images=:detail_images, status=:status
		WHERE id=:id`
	if _, err := r.db.NamedExecContext(ctx, query, p); err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	return nil
}

func (r *productRepository) Delete(ctx context.Context, id uint64) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id = $1`, id); err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	return nil
}
