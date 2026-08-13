package model

import "time"

// Document maps to the `documents` table.
// DeletedAt is nil for active documents; a non-nil value means the document
// has been soft-deleted and can still be brought back via Restore.
type Document struct {
	ID          uint64     `db:"id"`
	UserID      uint64     `db:"user_id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Author      string     `db:"author"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

// DocumentOutline maps to the `document_outlines` table.
// ParentID is nil for top-level nodes.
type DocumentOutline struct {
	ID         uint64    `db:"id"`
	DocumentID uint64    `db:"document_id"`
	ParentID   *uint64   `db:"parent_id"`
	Title      string    `db:"title"`
	SortOrder  int       `db:"sort_order"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

// OutlineContent maps to the `outline_contents` table.
// It has a 1:1 relationship with DocumentOutline, keyed by OutlineID.
type OutlineContent struct {
	OutlineID uint64    `db:"outline_id"`
	Content   string    `db:"content"`
	UpdatedAt time.Time `db:"updated_at"`
}
