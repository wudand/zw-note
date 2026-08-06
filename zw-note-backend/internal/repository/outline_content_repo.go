package repository

import (
	"context"
	"database/sql"
	"fmt"

	"zw-note-backend/internal/model"

	"github.com/jmoiron/sqlx"
)

// OutlineContentRepository defines the data-access contract for the outline_contents table.
type OutlineContentRepository interface {
	GetByOutlineID(ctx context.Context, outlineID uint64) (*model.OutlineContent, error)
	// Upsert supports an optional tx since it is typically called right after
	// creating the owning outline node within the same transaction.
	Upsert(ctx context.Context, tx *sqlx.Tx, outlineID uint64, content string) error
}

type outlineContentRepository struct {
	db *sqlx.DB
}

func NewOutlineContentRepository(db *sqlx.DB) OutlineContentRepository {
	return &outlineContentRepository{db: db}
}

func (r *outlineContentRepository) GetByOutlineID(ctx context.Context, outlineID uint64) (*model.OutlineContent, error) {
	c := &model.OutlineContent{}
	err := r.db.GetContext(ctx, c, `SELECT * FROM outline_contents WHERE outline_id = $1`, outlineID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get outline content by outline id: %w", err)
	}
	return c, nil
}

func (r *outlineContentRepository) Upsert(ctx context.Context, tx *sqlx.Tx, outlineID uint64, content string) error {
	query := `INSERT INTO outline_contents (outline_id, content)
		VALUES ($1, $2)
		ON CONFLICT (outline_id) DO UPDATE SET content = EXCLUDED.content`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, outlineID, content)
	} else {
		_, err = r.db.ExecContext(ctx, query, outlineID, content)
	}
	if err != nil {
		return fmt.Errorf("upsert outline content: %w", err)
	}
	return nil
}
