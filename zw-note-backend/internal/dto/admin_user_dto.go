package dto

// CreateAdminUserRequest is the payload for POST /api/admin/v1/users.
type CreateAdminUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Role     string `json:"role"     binding:"required,oneof=super_admin admin viewer"`
}

// UpdateAdminUserRequest is the payload for PUT /api/admin/v1/users/:id.
type UpdateAdminUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=50"`
	Role     string `json:"role"     binding:"omitempty,oneof=super_admin admin viewer"`
	Status   *int8  `json:"status"   binding:"omitempty,oneof=0 1"`
}

// AdminLoginRequest is the payload for POST /api/admin/v1/auth/login.
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminLoginResponse carries the admin JWT token and its expiry.
type AdminLoginResponse struct {
	Token    string `json:"token"`
	ExpireAt int64  `json:"expire_at"`
	Role     string `json:"role"`
}

// AdminUserResponse is the public-facing representation of an admin user.
type AdminUserResponse struct {
	ID        uint64 `json:"id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	Status    int8   `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// AdminUserListResponse wraps a paginated admin user list.
type AdminUserListResponse struct {
	Total int64                `json:"total"`
	List  []*AdminUserResponse `json:"list"`
}
