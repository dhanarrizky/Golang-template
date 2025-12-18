package middleware

import (
	"net/http"
	"strings"

	ports "github.com/dhanarrizky/Golang-template/internal/ports/auth"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates access token and places user id into gin.Context
func AuthMiddleware(tokenSigner ports.TokenSigner) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization format",
			})
			return
		}

		tokenStr := parts[1]

		payload, err := tokenSigner.VerifyAccessToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		// 🔥 SET CLAIMS KE CONTEXT
		c.Set("userID", payload.UserID)
		c.Set("tokenID", payload.TokenID)
		c.Set("tokenExpiresAt", payload.ExpiresAt)

		c.Next()
	}
}
