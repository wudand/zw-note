package model

import "time"

// AdminUser maps to the `admin_users` table (management console).
type AdminUser struct {
	ID        uint64    `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	Status    int8      `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Admin roles – stored in JWT claims for stateless RBAC.
const (
	RoleSuperAdmin = "super_admin" // full access
	RoleAdmin      = "admin"       // manage resources, cannot manage admins
	RoleViewer     = "viewer"      // read-only
)

// AdminRoleWeight returns a numeric weight for role comparison.
// Higher weight means broader permissions.
func AdminRoleWeight(role string) int {
	switch role {
	case RoleSuperAdmin:
		return 100
	case RoleAdmin:
		return 50
	case RoleViewer:
		return 10
	default:
		return 0
	}
}
