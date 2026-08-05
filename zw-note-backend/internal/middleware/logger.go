package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go-web-api/pkg/utils"
	"go.uber.org/zap"
)

// RequestLogger records method, path, status, latency and client IP for every request.
func RequestLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				log.Error("request error", append(fields, zap.String("error", e))...)
			}
			return
		}

		switch {
		case status >= http.StatusInternalServerError:
			log.Error("request", fields...)
		case status >= http.StatusBadRequest:
			log.Warn("request", fields...)
		default:
			log.Info("request", fields...)
		}
	}
}

// Recovery catches panics, logs them with zap, and responds with 500.
func Recovery(log *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, err interface{}) {
		log.Error("panic recovered",
			zap.Any("error", err),
			zap.String("path", c.Request.URL.Path),
		)
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Response{
			Code:    utils.CodeServerError,
			Message: "internal server error",
		})
	})
}
