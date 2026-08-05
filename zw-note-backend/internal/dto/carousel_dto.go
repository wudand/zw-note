package dto

// CreateCarouselRequest is the payload for creating a carousel.
type CreateCarouselRequest struct {
	ImageURL  string `json:"image_url"  binding:"required,max=500"`
	Title     string `json:"title"      binding:"max=100"`
	Link      string `json:"link"      binding:"max=500"`
	SortOrder int    `json:"sort_order" binding:"min=0"`
}

// UpdateCarouselRequest is the payload for updating a carousel.
type UpdateCarouselRequest struct {
	ImageURL  *string `json:"image_url"  binding:"omitempty,max=500"`
	Title     *string `json:"title"      binding:"omitempty,max=100"`
	Link      *string `json:"link"      binding:"omitempty,max=500"`
	SortOrder *int    `json:"sort_order" binding:"omitempty,min=0"`
	Status    *int8   `json:"status"    binding:"omitempty,oneof=0 1"`
}

// CarouselResponse is the public-facing representation of a carousel.
type CarouselResponse struct {
	ID        uint64 `json:"id"`
	ImageURL  string `json:"image_url"`
	Title     string `json:"title"`
	Link      string `json:"link"`
	SortOrder int    `json:"sort_order"`
	Status    int8   `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CarouselListResponse wraps a list of carousels.
type CarouselListResponse struct {
	List []*CarouselResponse `json:"list"`
}
