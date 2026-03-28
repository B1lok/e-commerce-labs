package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
				log.Printf("[ERROR] %s (wrapped: %v)", apiErr.Description, apiErr.Err)
			}
			c.AbortWithStatusJSON(apiErr.Code, gin.H{"description": apiErr.Description})
			return
		}
		log.Printf("[ERROR] unhandled error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"description": "Internal server error"})
	}
}
