package middleware

import (
	"net/http"
	"sync"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Simple in-memory rate limiter per key (IP or API key).
// Not distributed — for single-process use. For multi-instance, plug Redis or other store.
type simpleBucket struct {
	remaining int
	reset     time.Time
}

type rateLimiterStore struct {
	mu    sync.Mutex
	store map[string]*simpleBucket
	// configuration
	limit  int
	window time.Duration
}

func newRateLimiterStore(limit int, window time.Duration) *rateLimiterStore {
	return &rateLimiterStore{
		store:  make(map[string]*simpleBucket),
		limit:  limit,
		window: window,
	}
}

// RateLimitMiddleware returns a gin middleware. keyFunc decides the identity (IP, header, etc).
// limit = number of requests per window (e.g., 60 per minute).
func RateLimitMiddleware(limit int, window time.Duration, keyFunc func(c *gin.Context) string) gin.HandlerFunc {
	store := newRateLimiterStore(limit, window)

	return func(c *gin.Context) {
		key := keyFunc(c)
		if key == "" {
			key = c.ClientIP()
		}
		now := time.Now()

		store.mu.Lock()
		b, ok := store.store[key]
		if !ok || now.After(b.reset) {
			b = &simpleBucket{
				remaining: store.limit - 1,
				reset:     now.Add(store.window),
			}
			store.store[key] = b
			store.mu.Unlock()
		} else {
			if b.remaining <= 0 {
				// rate limited
				resetIn := int(b.reset.Sub(now).Seconds())
				store.mu.Unlock()
				c.Header("Retry-After", string(resetIn))
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"error":      "rate limit exceeded",
					"retry_after": resetIn,
				})
				return
			}
			b.remaining--
			store.mu.Unlock()
		}

		// set headers for visibility
		c.Header("X-RateLimit-Limit", itoa(store.limit))
		c.Header("X-RateLimit-Remaining", itoa(store.store[key].remaining))
		c.Header("X-RateLimit-Reset", store.store[key].reset.Format(time.RFC3339))

		c.Next()
	}
}

// helper simple itoa for small integers
func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
