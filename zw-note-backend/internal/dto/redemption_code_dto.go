package dto

// CreateRedemptionCodesRequest 批量创建兑换码（创建时绑定产品）
type CreateRedemptionCodesRequest struct {
	ProductIDs []uint64 `json:"product_ids" binding:"required,min=1,max=20,dive,min=1"`
	Count      int     `json:"count"       binding:"required,min=1,max=50"`
}

// UpdateRedemptionCodeRequest 更新兑换码（仅未使用的可更新绑定产品）
type UpdateRedemptionCodeRequest struct {
	ProductIDs []uint64 `json:"product_ids" binding:"required,min=1,max=20,dive,min=1"`
}

// RedemptionCodeListQuery 列表查询
type RedemptionCodeListQuery struct {
	Page     int   `form:"page"      binding:"omitempty,min=1"`
	PageSize int   `form:"page_size" binding:"omitempty,min=1,max=100"`
	Status   *int8 `form:"status"    binding:"omitempty,oneof=0 1"`
}

// RedemptionCodeResponse 兑换码响应
type RedemptionCodeResponse struct {
	ID             uint64   `json:"id"`
	Code           string   `json:"code"`
	Status         int8     `json:"status"`
	MPUserID       *uint64  `json:"mp_user_id,omitempty"`
	UsedProductID  *uint64  `json:"used_product_id,omitempty"`
	UsedAt         *string  `json:"used_at,omitempty"`
	ProductIDs     []uint64 `json:"product_ids"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

// RedemptionCodeListResponse 兑换码列表
type RedemptionCodeListResponse struct {
	Total int64                    `json:"total"`
	List  []*RedemptionCodeResponse `json:"list"`
}

// ValidateRedemptionCodeRequest 校验兑换码（MP 端）
type ValidateRedemptionCodeRequest struct {
	Code string `json:"code" binding:"required,len=12"`
}

// ValidateRedemptionCodeResponse 校验成功返回可选产品
type ValidateRedemptionCodeResponse struct {
	RedemptionCodeID uint64            `json:"redemption_code_id"`
	Products         []*ProductSummary `json:"products"`
}

// ProductSummary 产品摘要（用于兑换码可选产品列表）
type ProductSummary struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	CoverImage  string `json:"cover_image"`
	Specification string `json:"specification"`
}

// SelectProductRequest 选择产品（MP 端）
type SelectProductRequest struct {
	ProductID uint64 `json:"product_id" binding:"required,min=1"`
}
