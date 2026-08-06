package model

import (
	"time"

	"zw-note-backend/pkg/types"
)

// Product maps to the `products` table.
type Product struct {
	ID             uint64           `db:"id"`
	CategoryID     uint64           `db:"category_id"`
	Name           string          `db:"name"`
	Ingredients    string           `db:"ingredients"`
	Grade          string           `db:"grade"`
	Storage        string          `db:"storage"`
	Specification  string          `db:"specification"`
	CoverImage     string          `db:"cover_image"`
	CarouselImages types.StringSlice `db:"carousel_images"`
	Detail         string          `db:"detail"`
	DetailImages   types.StringSlice `db:"detail_images"`
	Status         int8            `db:"status"`
	CreatedAt      time.Time       `db:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at"`
}
