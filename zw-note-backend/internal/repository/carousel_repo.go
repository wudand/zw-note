package repository

import (
	"context"
	"database/sql"
	"fmt"

	"zw-note-backend/internal/model"

	"github.com/jmoiron/sqlx"
)

// CarouselRepository defines the data-access contract for the carousels table.
type CarouselRepository interface {
	Create(ctx context.Context, c *model.Carousel) error
	GetByID(ctx context.Context, id uint64) (*model.Carousel, error)
	List(ctx context.Context, enabledOnly bool) ([]*model.Carousel, error)
	CountEnabled(ctx context.Context) (int, error)
	Update(ctx context.Context, c *model.Carousel) error
	Delete(ctx context.Context, id uint64) error
}

type carouselRepository struct {
	db *sqlx.DB
}

func NewCarouselRepository(db *sqlx.DB) CarouselRepository {
	return &carouselRepository{db: db}
}

func (r *carouselRepository) Create(ctx context.Context, c *model.Carousel) error {
	query := `INSERT INTO carousels (image_url, title, link, sort_order, status)
		VALUES (:image_url, :title, :link, :sort_order, :status) RETURNING id`
	if err := namedReturningID(ctx, r.db, query, c, &c.ID); err != nil {
		return fmt.Errorf("create carousel: %w", err)
	}
	return nil
}

func (r *carouselRepository) GetByID(ctx context.Context, id uint64) (*model.Carousel, error) {
	c := &model.Carousel{}
	err := r.db.GetContext(ctx, c, `SELECT * FROM carousels WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get carousel by id: %w", err)
	}
	return c, nil
}

func (r *carouselRepository) List(ctx context.Context, enabledOnly bool) ([]*model.Carousel, error) {
	query := `SELECT * FROM carousels`
	if enabledOnly {
		query += ` WHERE status = 1`
	}
	query += ` ORDER BY sort_order ASC, id ASC`

	var list []*model.Carousel
	if err := r.db.SelectContext(ctx, &list, query); err != nil {
		return nil, fmt.Errorf("list carousels: %w", err)
	}
	return list, nil
}

func (r *carouselRepository) CountEnabled(ctx context.Context) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM carousels WHERE status = 1`)
	if err != nil {
		return 0, fmt.Errorf("count enabled carousels: %w", err)
	}
	return count, nil
}

func (r *carouselRepository) Update(ctx context.Context, c *model.Carousel) error {
	query := `UPDATE carousels SET image_url=:image_url, title=:title, link=:link, sort_order=:sort_order, status=:status WHERE id=:id`
	if _, err := r.db.NamedExecContext(ctx, query, c); err != nil {
		return fmt.Errorf("update carousel: %w", err)
	}
	return nil
}

func (r *carouselRepository) Delete(ctx context.Context, id uint64) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM carousels WHERE id = $1`, id); err != nil {
		return fmt.Errorf("delete carousel: %w", err)
	}
	return nil
}
