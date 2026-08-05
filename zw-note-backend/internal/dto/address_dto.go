package dto

// CreateAddressRequest is the payload for creating an address.
type CreateAddressRequest struct {
	Receiver string `json:"receiver" binding:"required,max=50"`
	Phone    string `json:"phone"    binding:"required,max=20"`
	Province string `json:"province" binding:"required,max=50"`
	City     string `json:"city"     binding:"required,max=50"`
	District string `json:"district" binding:"required,max=50"`
	Detail   string `json:"detail"   binding:"required,max=255"`
	Tag      string `json:"tag"      binding:"omitempty,max=20"`
}

// UpdateAddressRequest is the payload for updating an address.
type UpdateAddressRequest struct {
	Receiver *string `json:"receiver" binding:"omitempty,max=50"`
	Phone    *string `json:"phone"    binding:"omitempty,max=20"`
	Province *string `json:"province" binding:"omitempty,max=50"`
	City     *string `json:"city"     binding:"omitempty,max=50"`
	District *string `json:"district" binding:"omitempty,max=50"`
	Detail   *string `json:"detail"   binding:"omitempty,max=255"`
	Tag      *string `json:"tag"      binding:"omitempty,max=20"`
}

// AddressResponse is the public-facing representation of an address.
type AddressResponse struct {
	ID        uint64 `json:"id"`
	MPUserID  uint64 `json:"mp_user_id"`
	Receiver  string `json:"receiver"`
	Phone     string `json:"phone"`
	Province  string `json:"province"`
	City      string `json:"city"`
	District  string `json:"district"`
	Detail    string `json:"detail"`
	Tag       string `json:"tag"`
	IsDefault bool   `json:"is_default"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// AddressListResponse wraps a list of addresses.
type AddressListResponse struct {
	List []*AddressResponse `json:"list"`
}
