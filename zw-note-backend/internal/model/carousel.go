package model

import "time"

// Carousel maps to the `carousels` table.
type Carousel struct {
	ID        uint64    `db:"id"`
	ImageURL  string    `db:"image_url"`
	Title     string    `db:"title"`
	Link      string    `db:"link"`
	SortOrder int       `db:"sort_order"`
	Status    int8      `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

const (
	CarouselStatusEnabled  int8 = 1
	CarouselStatusDisabled int8 = 0
)
