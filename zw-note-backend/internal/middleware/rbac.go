package middleware

import (
	"zw-note-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// RequireRole returns a middleware that aborts with 403 if the authenticated
// admin's role is not among the allowed roles.
//
// Must be placed after AdminJWTAuth in the middleware chain:
//
//	group.DELETE("/users/:id",
//	    middleware.AdminJWTAuth(cfg.JWT),
//	    middleware.RequireRole("super_admin"),
//	    handler.DeleteUser,
//	)
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(c *gin.Context) {
		role, exists := c.Get(CtxKeyRole)
		if !exists {
			utils.Forbidden(c)
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			utils.Forbidden(c)
			c.Abort()
			return
		}

		if _, permitted := allowed[roleStr]; !permitted {
			utils.Forbidden(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
