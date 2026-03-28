package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
)

func TransactionMiddleware(txManager TxManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.IsAborted() {
			return
		}
		_ = txManager.ExecuteTx(c.Request.Context(), func(txCtx context.Context) error {
			c.Request = c.Request.WithContext(txCtx)

			c.Next()

			if len(c.Errors) > 0 {
				return c.Errors[0]
			}

			return nil
		})
	}
}
