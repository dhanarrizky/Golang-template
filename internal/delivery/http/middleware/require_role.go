package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {

		rolesAny, exists := c.Get("roles")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "roles not found in token",
			})
			return
		}

		roles, ok := rolesAny.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "invalid roles format",
			})
			return
		}

		for _, role := range roles {
			if role == requiredRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "insufficient permissions",
		})
	}
}

func RequireAnyRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		rolesAny, _ := c.Get("roles")
		roles := rolesAny.([]string)

		for _, r := range roles {
			for _, req := range requiredRoles {
				if r == req {
					c.Next()
					return
				}
			}
		}

		c.AbortWithStatusJSON(403, gin.H{
			"error": "insufficient permissions",
		})
	}
}
