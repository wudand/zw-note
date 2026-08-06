package repository

import (
	"context"
	"database/sql"
	"fmt"

	"zw-note-backend/internal/model"

	"github.com/jmoiron/sqlx"
)

// OutlineRepository defines the data-access contract for the document_outlines table.
type OutlineRepository interface {
	// Create supports an optional tx since a new node is typically created together
	// with its (initially empty) content row.
	Create(ctx context.Context, tx *sqlx.Tx, o *model.DocumentOutline) error
	GetByID(ctx context.Context, id uint64) (*model.DocumentOutline, error)
	// ListByDocument returns all nodes of a document as a flat slice; the service
	// layer assembles them into a tree.
	ListByDocument(ctx context.Context, documentID uint64) ([]*model.DocumentOutline, error)
	Update(ctx context.Context, tx *sqlx.Tx, o *model.DocumentOutline) error
	// Delete relies on the FK's ON DELETE CASCADE to remove descendant nodes and content.
	Delete(ctx context.Context, tx *sqlx.Tx, id uint64) error
	// IsDescendant reports whether candidateAncestorID lies within the subtree rooted
	// at nodeID (i.e. is a descendant of nodeID) inside the given document. It is used
	// to reject moves that would turn a node into a descendant of itself (a cycle).
	IsDescendant(ctx context.Context, documentID, nodeID, candidateAncestorID uint64) (bool, error)
}

type outlineRepository struct {
	db *sqlx.DB
}

func NewOutlineRepository(db *sqlx.DB) OutlineRepository {
	return &outlineRepository{db: db}
}

func (r *outlineRepository) Create(ctx context.Context, tx *sqlx.Tx, o *model.DocumentOutline) error {
	query := `INSERT INTO document_outlines (document_id, parent_id, title, sort_order)
		VALUES (:document_id, :parent_id, :title, :sort_order)
		RETURNING id`
	var q sqlx.QueryerContext = r.db
	if tx != nil {
		q = tx
	}
	if err := namedReturningID(ctx, q, query, o, &o.ID); err != nil {
		return fmt.Errorf("create outline: %w", err)
	}
	return nil
}

func (r *outlineRepository) GetByID(ctx context.Context, id uint64) (*model.DocumentOutline, error) {
	o := &model.DocumentOutline{}
	err := r.db.GetContext(ctx, o, `SELECT * FROM document_outlines WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get outline by id: %w", err)
	}
	return o, nil
}

func (r *outlineRepository) ListByDocument(ctx context.Context, documentID uint64) ([]*model.DocumentOutline, error) {
	var list []*model.DocumentOutline
	err := r.db.SelectContext(ctx, &list,
		`SELECT * FROM document_outlines WHERE document_id = $1 ORDER BY parent_id NULLS FIRST, sort_order`,
		documentID)
	if err != nil {
		return nil, fmt.Errorf("list outlines by document: %w", err)
	}
	return list, nil
}

func (r *outlineRepository) Update(ctx context.Context, tx *sqlx.Tx, o *model.DocumentOutline) error {
	query := `UPDATE document_outlines SET parent_id=:parent_id, title=:title, sort_order=:sort_order
		WHERE id=:id`
	if tx != nil {
		if _, err := tx.NamedExecContext(ctx, query, o); err != nil {
			return fmt.Errorf("update outline: %w", err)
		}
	} else {
		if _, err := r.db.NamedExecContext(ctx, query, o); err != nil {
			return fmt.Errorf("update outline: %w", err)
		}
	}
	return nil
}

func (r *outlineRepository) Delete(ctx context.Context, tx *sqlx.Tx, id uint64) error {
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, `DELETE FROM document_outlines WHERE id = $1`, id)
	} else {
		_, err = r.db.ExecContext(ctx, `DELETE FROM document_outlines WHERE id = $1`, id)
	}
	if err != nil {
		return fmt.Errorf("delete outline: %w", err)
	}
	return nil
}

func (r *outlineRepository) IsDescendant(ctx context.Context, documentID, nodeID, candidateAncestorID uint64) (bool, error) {
	query := `WITH RECURSIVE descendants AS (
		SELECT id FROM document_outlines WHERE parent_id = $1 AND document_id = $2
		UNION ALL
		SELECT o.id FROM document_outlines o
		INNER JOIN descendants d ON o.parent_id = d.id
		WHERE o.document_id = $2
	)
	SELECT EXISTS (SELECT 1 FROM descendants WHERE id = $3)`
	var exists bool
	if err := r.db.GetContext(ctx, &exists, query, nodeID, documentID, candidateAncestorID); err != nil {
		return false, fmt.Errorf("check outline descendant: %w", err)
	}
	return exists, nil
}
