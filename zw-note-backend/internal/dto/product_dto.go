package dto

// CreateProductRequest is the payload for POST /api/admin/v1/products.
type CreateProductRequest struct {
	CategoryID     uint64   `json:"category_id"    binding:"required"`
	Name           string   `json:"name"           binding:"required,min=1,max=100"`
	Ingredients    string   `json:"ingredients"    binding:"max=255"`
	Grade          string   `json:"grade"          binding:"max=50"`
	Storage        string   `json:"storage"        binding:"max=255"`
	Specification  string   `json:"specification"  binding:"max=100"`
	CoverImage     string   `json:"cover_image"    binding:"required,max=500"`
	CarouselImages []string `json:"carousel_images" binding:"omitempty"`
	Detail         string   `json:"detail"         binding:"omitempty"`
	DetailImages   []string `json:"detail_images"  binding:"omitempty"`
}

// UpdateProductRequest is the payload for PUT /api/admin/v1/products/:id.
// Pointer fields: nil = no change, non-nil = update (including empty).
type UpdateProductRequest struct {
	CategoryID     *uint64  `json:"category_id"     binding:"omitempty"`
	Name           *string  `json:"name"           binding:"omitempty,max=100"`
	Ingredients    *string  `json:"ingredients"    binding:"omitempty,max=255"`
	Grade          *string  `json:"grade"          binding:"omitempty,max=50"`
	Storage        *string  `json:"storage"        binding:"omitempty,max=255"`
	Specification  *string  `json:"specification"  binding:"omitempty,max=100"`
	CoverImage     *string  `json:"cover_image"    binding:"omitempty,max=500"`
	CarouselImages []string `json:"carousel_images" binding:"omitempty"`
	Detail         *string  `json:"detail"         binding:"omitempty"`
	DetailImages   []string `json:"detail_images"  binding:"omitempty"`
	Status         *int8    `json:"status"        binding:"omitempty,oneof=0 1"`
}

// ProductListQuery holds list filter params.
type ProductListQuery struct {
	Page       int     `form:"page"        binding:"omitempty,min=1"`
	PageSize   int     `form:"page_size"   binding:"omitempty,min=1,max=100"`
	CategoryID *uint64 `form:"category_id" binding:"omitempty"`
}

// ProductResponse is the public-facing representation of a product.
type ProductResponse struct {
	ID             uint64   `json:"id"`
	CategoryID     uint64   `json:"category_id"`
	Name           string   `json:"name"`
	Ingredients    string   `json:"ingredients"`
	Grade          string   `json:"grade"`
	Storage        string   `json:"storage"`
	Specification  string   `json:"specification"`
	CoverImage     string   `json:"cover_image"`
	CarouselImages []string `json:"carousel_images"`
	Detail         string   `json:"detail"`
	DetailImages   []string `json:"detail_images"`
	Status         int8     `json:"status"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

// ProductListResponse wraps a paginated product list.
type ProductListResponse struct {
	Total int64              `json:"total"`
	List  []*ProductResponse `json:"list"`
}
