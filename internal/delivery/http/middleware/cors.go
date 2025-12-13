package middleware

import "github.com/gin-gonic/gin"

// Simple CORS middleware. Use for APIs that need to be opened to browsers.
// You can enhance to read allowed origins from configuration.
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	// convert to map for quick lookup if provided
	origins := map[string]struct{}{}
	for _, o := range allowedOrigins {
		origins[o] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		// If allowedOrigins is empty, allow any origin (use with caution)
		allowAny := len(origins) == 0

		if allowAny {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			if _, ok := origins[origin]; ok {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Requested-With, X-Api-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		// Preflight
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
