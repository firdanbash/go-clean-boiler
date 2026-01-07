package middleware

import (
	"net/http"

	"github.com/firdanbash/go-clean-boiler/pkg/logger"
	"github.com/firdanbash/go-clean-boiler/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorMiddleware handles panics and errors
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
				)

				response.InternalServerError(c, "Internal server error", nil)
				c.Abort()
			}
		}()

		c.Next()

		// Check if there were any errors during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			logger.Error("Request error",
				zap.Error(err.Err),
				zap.String("path", c.Request.URL.Path),
			)

			// If response hasn't been written yet
			if c.Writer.Status() == http.StatusOK {
				response.InternalServerError(c, "An error occurred", err.Error())
			}
		}
	}
}
