package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// TransactionManager provides transaction management. Services depend on this
// interface instead of *sqlx.DB for clearer semantics: "transaction capability"
// rather than "database connection".
type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(tx *sqlx.Tx) error) error
}

type transactionManager struct {
	db *sqlx.DB
}

// NewTransactionManager creates a TransactionManager backed by the given DB.
func NewTransactionManager(db *sqlx.DB) TransactionManager {
	return &transactionManager{db: db}
}

// WithTransaction runs fn within a transaction. On success it commits; on error
// or panic it rolls back. defer tx.Rollback() is safe after Commit (no-op).
func (m *transactionManager) WithTransaction(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}
