package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"zw-note-backend/internal/config"
)

// CORS returns a middleware that sets Access-Control-* headers based on cfg.
// Preflight OPTIONS requests are terminated with 204 No Content.
func CORS(cfg config.CORSConfig) gin.HandlerFunc {
	allowMethods := strings.Join(cfg.AllowMethods, ", ")
	allowHeaders := strings.Join(cfg.AllowHeaders, ", ")
	maxAge := strconv.Itoa(cfg.MaxAge)

	originSet := make(map[string]struct{}, len(cfg.AllowOrigins))
	wildcardAll := false
	for _, o := range cfg.AllowOrigins {
		if o == "*" {
			wildcardAll = true
		}
		originSet[o] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		allowed := wildcardAll
		if !allowed {
			_, allowed = originSet[origin]
		}

		if allowed {
			out := origin
			if out == "" {
				out = "*"
			}
			c.Header("Access-Control-Allow-Origin", out)
			c.Header("Access-Control-Allow-Methods", allowMethods)
			c.Header("Access-Control-Allow-Headers", allowHeaders)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", maxAge)
			c.Header("Vary", "Origin")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
