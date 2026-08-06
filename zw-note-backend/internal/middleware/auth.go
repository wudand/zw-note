package middleware

import (
	"strings"

	"zw-note-backend/internal/config"
	"zw-note-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Context keys injected by auth middleware for downstream handlers.
const (
	CtxKeyUserID   = "user_id"
	CtxKeyUsername = "username"
	CtxKeyRole     = "role"
	CtxKeyOpenID   = "openid"
	CtxKeyAppID    = "app_id"
)

// AdminJWTAuth validates the Bearer token for the management console.
// On success it stores user_id, username, and role in the gin context.
func AdminJWTAuth(cfg config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, ok := extractBearerToken(c)
		if !ok {
			utils.Unauthorized(c)
			c.Abort()
			return
		}

		claims, err := utils.ParseAdminToken(tokenStr, cfg.AdminSecret)
		if err != nil {
			utils.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set(CtxKeyUserID, claims.UserID)
		c.Set(CtxKeyUsername, claims.Username)
		c.Set(CtxKeyRole, claims.Role)
		c.Next()
	}
}

// MPJWTAuth validates the Bearer token for WeChat mini-program users.
// On success it stores user_id, openid, and app_id in the gin context.
func MPJWTAuth(cfg config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, ok := extractBearerToken(c)
		if !ok {
			utils.Unauthorized(c)
			c.Abort()
			return
		}

		claims, err := utils.ParseMPToken(tokenStr, cfg.MPSecret)
		if err != nil {
			utils.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set(CtxKeyUserID, claims.UserID)
		c.Set(CtxKeyOpenID, claims.OpenID)
		c.Set(CtxKeyAppID, claims.AppID)
		c.Next()
	}
}

func extractBearerToken(c *gin.Context) (string, bool) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return "", false
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}
