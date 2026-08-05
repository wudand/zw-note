package model

import "time"

// Category maps to the `product_categories` table.
type Category struct {
	ID        uint64    `db:"id"`
	Name      string    `db:"name"`
	SortOrder int       `db:"sort_order"`
	Status    int8      `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
