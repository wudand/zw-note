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
	// GetByIDAndUser only returns active (non-deleted) documents.
	GetByIDAndUser(ctx context.Context, id, userID uint64) (*model.Document, error)
	// GetDeletedByIDAndUser only returns soft-deleted documents, for Restore.
	GetDeletedByIDAndUser(ctx context.Context, id, userID uint64) (*model.Document, error)
	// ListByUser only returns active (non-deleted) documents.
	ListByUser(ctx context.Context, userID uint64) ([]*model.Document, error)
	// ListDeletedByUser returns soft-deleted documents, most recently deleted first.
	ListDeletedByUser(ctx context.Context, userID uint64) ([]*model.Document, error)
	Update(ctx context.Context, tx *sqlx.Tx, d *model.Document) error
	// SoftDelete marks a document as deleted by setting deleted_at; it does not
	// physically remove the row (or its outlines/contents), so Restore can undo it.
	SoftDelete(ctx context.Context, tx *sqlx.Tx, id, userID uint64) error
	// Restore clears deleted_at, bringing a soft-deleted document back.
	Restore(ctx context.Context, tx *sqlx.Tx, id, userID uint64) error
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
	err := r.db.GetContext(ctx, d,
		`SELECT * FROM documents WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`, id, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get document by id and user: %w", err)
	}
	return d, nil
}

func (r *documentRepository) GetDeletedByIDAndUser(ctx context.Context, id, userID uint64) (*model.Document, error) {
	d := &model.Document{}
	err := r.db.GetContext(ctx, d,
		`SELECT * FROM documents WHERE id = $1 AND user_id = $2 AND deleted_at IS NOT NULL`, id, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get deleted document by id and user: %w", err)
	}
	return d, nil
}

func (r *documentRepository) ListByUser(ctx context.Context, userID uint64) ([]*model.Document, error) {
	var list []*model.Document
	err := r.db.SelectContext(ctx, &list,
		`SELECT * FROM documents WHERE user_id = $1 AND deleted_at IS NULL ORDER BY updated_at DESC`,
		userID)
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	return list, nil
}

func (r *documentRepository) ListDeletedByUser(ctx context.Context, userID uint64) ([]*model.Document, error) {
	var list []*model.Document
	err := r.db.SelectContext(ctx, &list,
		`SELECT * FROM documents WHERE user_id = $1 AND deleted_at IS NOT NULL ORDER BY deleted_at DESC`,
		userID)
	if err != nil {
		return nil, fmt.Errorf("list deleted documents: %w", err)
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

func (r *documentRepository) SoftDelete(ctx context.Context, tx *sqlx.Tx, id, userID uint64) error {
	query := `UPDATE documents SET deleted_at = NOW()
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, id, userID)
	} else {
		_, err = r.db.ExecContext(ctx, query, id, userID)
	}
	if err != nil {
		return fmt.Errorf("soft delete document: %w", err)
	}
	return nil
}

func (r *documentRepository) Restore(ctx context.Context, tx *sqlx.Tx, id, userID uint64) error {
	query := `UPDATE documents SET deleted_at = NULL
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NOT NULL`
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, id, userID)
	} else {
		_, err = r.db.ExecContext(ctx, query, id, userID)
	}
	if err != nil {
		return fmt.Errorf("restore document: %w", err)
	}
	return nil
}
