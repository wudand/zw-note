package dto

// CreateCouponRequest is the payload for creating a coupon.
type CreateCouponRequest struct {
	Name          string `json:"name"           binding:"required,max=100"`
	Type          string `json:"type"           binding:"required,oneof=new_user spend_reduce"`
	MinAmount     int    `json:"min_amount"     binding:"required,min=0"` // 分
	DiscountValue int    `json:"discount_value" binding:"required,min=1"` // 分
	ValidDays     int    `json:"valid_days"    binding:"required,min=1,max=365"`
}

// UpdateCouponRequest is the payload for updating a coupon.
type UpdateCouponRequest struct {
	Name          *string `json:"name"           binding:"omitempty,max=100"`
	MinAmount     *int    `json:"min_amount"     binding:"omitempty,min=0"`
	DiscountValue *int    `json:"discount_value" binding:"omitempty,min=1"`
	ValidDays     *int    `json:"valid_days"    binding:"omitempty,min=1,max=365"`
	Status        *int8   `json:"status"        binding:"omitempty,oneof=0 1"`
}

// CouponListQuery holds list filter params.
type CouponListQuery struct {
	Page     int     `form:"page"      binding:"omitempty,min=1"`
	PageSize int     `form:"page_size" binding:"omitempty,min=1,max=100"`
	Type     *string `form:"type"      binding:"omitempty,oneof=new_user spend_reduce"`
	Status   *int8   `form:"status"   binding:"omitempty,oneof=0 1"`
}

// CouponResponse is the public-facing representation of a coupon template.
type CouponResponse struct {
	ID            uint64 `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	MinAmount     int    `json:"min_amount"`     // 分
	DiscountValue int    `json:"discount_value"` // 分
	ValidDays     int    `json:"valid_days"`
	Stackable     bool   `json:"stackable"`
	Status        int8   `json:"status"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// CouponListResponse wraps a paginated coupon list.
type CouponListResponse struct {
	Total int64             `json:"total"`
	List  []*CouponResponse `json:"list"`
}

// UserCouponResponse is the representation of a user's claimed coupon.
type UserCouponResponse struct {
	ID            uint64  `json:"id"`
	CouponID      uint64  `json:"coupon_id"`
	CouponName    string  `json:"coupon_name"`
	CouponType    string  `json:"coupon_type"`
	MinAmount     int     `json:"min_amount"`
	DiscountValue int     `json:"discount_value"`
	Stackable     bool    `json:"stackable"`
	Status        string  `json:"status"`
	ClaimedAt     string  `json:"claimed_at"`
	ExpiryAt      string  `json:"expiry_at"`
	UsedAt        *string `json:"used_at,omitempty"`
}

// UserCouponListResponse wraps a list of user coupons.
type UserCouponListResponse struct {
	List []*UserCouponResponse `json:"list"`
}
