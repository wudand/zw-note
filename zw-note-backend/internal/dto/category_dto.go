package dto

// CreateCategoryRequest is the payload for POST /api/admin/v1/categories.
type CreateCategoryRequest struct {
	Name      string `json:"name"      binding:"required,min=1,max=50"`
	SortOrder int    `json:"sort_order" binding:"omitempty"`
}

// UpdateCategoryRequest is the payload for PUT /api/admin/v1/categories/:id.
type UpdateCategoryRequest struct {
	Name      string `json:"name"      binding:"omitempty,min=1,max=50"`
	SortOrder *int   `json:"sort_order" binding:"omitempty"`
	Status    *int8  `json:"status"    binding:"omitempty,oneof=0 1"`
}

// CategoryResponse is the public-facing representation of a category.
type CategoryResponse struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
	Status    int8   `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CategoryListResponse wraps a paginated category list.
type CategoryListResponse struct {
	Total int64               `json:"total"`
	List  []*CategoryResponse `json:"list"`
}
