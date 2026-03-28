package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"weather-api/pkg/errors"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		errs := c.Errors
		if len(errs) == 0 {
			return
		}

		err := errs[0].Err

		if apiErr, ok := errors.IsAPIError(err); ok {
			if apiErr.Code >= 500 {
				slog.Error("API error", "description", apiErr.Description, "error", apiErr.Err)
			}
			c.AbortWithStatusJSON(apiErr.Code, gin.H{"description": apiErr.Description})
			return
		}
		slog.Error("Unhandled error", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"description": "Internal server error"})
	}
}
