package middleware

import "github.com/gin-gonic/gin"

// NotesAuth is a placeholder auth middleware for the notes module MVP, which
// does not yet have real login. It injects the configured default user ID
// into the gin context so handlers can read it the same way JWT-backed
// handlers do (see CtxKeyUserID). Once real authentication is added, this is
// the only place that needs to be swapped out.
func NotesAuth(defaultUserID uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(CtxKeyUserID, defaultUserID)
		c.Next()
	}
}
