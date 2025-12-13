package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// RecoveryWithLogging is a custom recovery that returns 500 and logs stack trace.
func RecoveryWithLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// log stack trace (stdout). Replace with structured logger as needed.
				fmt.Printf("PANIC recovered: %v\n%s\n", r, debug.Stack())
				// You can add more context (request id, path, headers) if desired.
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
			}
		}()
		c.Next()
	}
}
