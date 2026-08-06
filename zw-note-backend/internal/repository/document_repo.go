package repository

import (
	"context"
	"database/sql"
	"fmt"

	"zw-note-backend/internal/model"

	"github.com/jmoiron/sqlx"
)

// DocumentRepository defines the data-access contract for the documents table.
type DocumentRepository interface {
	// Create supports an optional tx since document creation is typically bundled
	// with creating a default root outline node and its empty content.
	Create(ctx context.Context, tx *sqlx.Tx, d *model.Document) error
	GetByIDAndUser(ctx context.Context, id, userID uint64) (*model.Document, error)
	ListByUser(ctx context.Context, userID uint64) ([]*model.Document, error)
	Update(ctx context.Context, tx *sqlx.Tx, d *model.Document) error
	Delete(ctx context.Context, tx *sqlx.Tx, id, userID uint64) error
}

type documentRepository struct {
	db *sqlx.DB
}

func NewDocumentRepository(db *sqlx.DB) DocumentRepository {
	return &documentRepository{db: db}
}

func (r *documentRepository) Create(ctx context.Context, tx *sqlx.Tx, d *model.Document) error {
	query := `INSERT INTO documents (user_id, title, description, author)
		VALUES (:user_id, :title, :description, :author)
		RETURNING id`
	var q sqlx.QueryerContext = r.db
	if tx != nil {
		q = tx
	}
	if err := namedReturningID(ctx, q, query, d, &d.ID); err != nil {
		return fmt.Errorf("create document: %w", err)
	}
	return nil
}

func (r *documentRepository) GetByIDAndUser(ctx context.Context, id, userID uint64) (*model.Document, error) {
	d := &model.Document{}
	err := r.db.GetContext(ctx, d, `SELECT * FROM documents WHERE id = $1 AND user_id = $2`, id, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get document by id and user: %w", err)
	}
	return d, nil
}

func (r *documentRepository) ListByUser(ctx context.Context, userID uint64) ([]*model.Document, error) {
	var list []*model.Document
	err := r.db.SelectContext(ctx, &list,
		`SELECT * FROM documents WHERE user_id = $1 ORDER BY updated_at DESC`,
		userID)
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	return list, nil
}

func (r *documentRepository) Update(ctx context.Context, tx *sqlx.Tx, d *model.Document) error {
	query := `UPDATE documents SET title=:title, description=:description, author=:author
		WHERE id=:id AND user_id=:user_id`
	if tx != nil {
		if _, err := tx.NamedExecContext(ctx, query, d); err != nil {
			return fmt.Errorf("update document: %w", err)
		}
	} else {
		if _, err := r.db.NamedExecContext(ctx, query, d); err != nil {
			return fmt.Errorf("update document: %w", err)
		}
	}
	return nil
}

func (r *documentRepository) Delete(ctx context.Context, tx *sqlx.Tx, id, userID uint64) error {
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, `DELETE FROM documents WHERE id = $1 AND user_id = $2`, id, userID)
	} else {
		_, err = r.db.ExecContext(ctx, `DELETE FROM documents WHERE id = $1 AND user_id = $2`, id, userID)
	}
	if err != nil {
		return fmt.Errorf("delete document: %w", err)
	}
	return nil
}
