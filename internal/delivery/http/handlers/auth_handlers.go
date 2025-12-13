package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/dhanarrizky/Golang-template/internal/domain/entities"
	"github.com/dhanarrizky/Golang-template/internal/config"
	"github.com/dhanarrizky/Golang-template/pkg/utils"
	"github.com/dhanarrizky/Golang-template/pkg/auth"
)

// Cookie name for refresh token
const RefreshCookieName = "refresh_token"

// AuthHandlers holds dependencies
type AuthHandlers struct {
	Signer      *auth.Signer
	RedisClient *redis.Client // may be nil
	Config      *config.Config
	UserSvc     UserAuthService
}

// UserAuthService minimal interface user usecase must implement
type UserAuthService interface {
	Authenticate(ctx context.Context, email, password string) (userID string, err error)
	GetByID(ctx context.Context, id string) (*entities.User, error)
	// optionally CreateUser etc.
}

// Login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginHandler authenticates user, issues access + refresh tokens.
// Refresh token is set as secure HttpOnly cookie.
func (h *AuthHandlers) LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body LoginRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}

		userID, err := h.UserSvc.Authenticate(c.Request.Context(), body.Email, body.Password)
		if err != nil {
			// avoid leaking whether email exists
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		// Create tokens
		access, _, err := h.Signer.NewAccessToken(userID, "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create access token"})
			return
		}

		// refresh token JTI random id (UUID recommended)
		jti := utils.GenerateID() // implement a secure random UUID v4 generator
		refreshToken, refreshClaims, err := h.Signer.NewRefreshToken(userID, jti)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create refresh token"})
			return
		}

		// Store jti in Redis for validation & rotation
		if h.RedisClient != nil {
			ttl := time.Until(refreshClaims.ExpiresAt.Time)
			key := "refresh:" + jti
			if err := h.RedisClient.Set(c.Request.Context(), key, userID, ttl).Err(); err != nil {
				// If Redis fails, treat based on policy: we can continue but log
				// For stricter security, return error. Here we'll fail-fast in production.
				if h.Config.IsProduction() {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
					return
				}
			}
		}

		// Set refresh token cookie
		httpOnly := true
		secure := h.Config.IsProduction()
		// sameSite := http.SameSiteLaxMode
		expSeconds := int(time.Until(refreshClaims.ExpiresAt.Time).Seconds())
		c.SetCookie(
			RefreshCookieName,
			refreshToken,
			expSeconds,
			"/",
			"",    // domain default to current
			secure,
			httpOnly,
		)
		// For SameSite you cannot set via SetCookie; use header if you need; Gin's c.SetCookie doesn't set SameSite prior to Go1.20, so you might set cookie manually via http.SetCookie.

		c.JSON(http.StatusOK, gin.H{
			"access_token": access,
			"token_type":   "bearer",
			"expires_in":   int(h.Signer.AccessTTL.Seconds()),
		})
	}
}


func (h *AuthHandlers) RefreshHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		rtCookie, err := c.Cookie(RefreshCookieName)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token required"})
			return
		}

		// Parse refresh token
		claims, err := h.Signer.ParseToken(rtCookie)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}
		jti := claims.ID
		if jti == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}

		ctx := c.Request.Context()

		// If Redis is enabled, validate jti exists and belongs to the same user
		if h.RedisClient != nil {
			key := "refresh:" + jti
			storedUser, err := h.RedisClient.Get(ctx, key).Result()
			if err == redis.Nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token revoked"})
				return
			} else if err != nil {
				// redis error
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
				return
			}
			if storedUser != claims.UserID {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token user mismatch"})
				return
			}
			// rotate: delete old key now to avoid reuse
			if err := h.RedisClient.Del(ctx, key).Err(); err != nil {
				// log but continue
			}
		}

		// Issue new access token
		access, _, err := h.Signer.NewAccessToken(claims.UserID, "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create access token"})
			return
		}

		// Create new refresh token with new jti
		newJTI := utils.GenerateID()
		refreshToken, refreshClaims, err := h.Signer.NewRefreshToken(claims.UserID, newJTI)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create refresh token"})
			return
		}

		// Store new jti in Redis
		if h.RedisClient != nil {
			ttl := time.Until(refreshClaims.ExpiresAt.Time)
			key := "refresh:" + newJTI
			if err := h.RedisClient.Set(ctx, key, claims.UserID, ttl).Err(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
				return
			}
		}

		// Set refresh cookie (rotate)
		secure := h.Config.IsProduction()
		httpOnly := true
		expSeconds := int(time.Until(refreshClaims.ExpiresAt.Time).Seconds())
		c.SetCookie(RefreshCookieName, refreshToken, expSeconds, "/", "", secure, httpOnly)

		c.JSON(http.StatusOK, gin.H{
			"access_token": access,
			"token_type":   "bearer",
			"expires_in":   int(h.Signer.AccessTTL.Seconds()),
		})
	}
}

func (h *AuthHandlers) LogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		rtCookie, err := c.Cookie(RefreshCookieName)
		if err == nil {
			// parse token to get jti
			if claims, err := h.Signer.ParseToken(rtCookie); err == nil {
				if h.RedisClient != nil && claims.ID != "" {
					key := "refresh:" + claims.ID
					_ = h.RedisClient.Del(c.Request.Context(), key).Err()
				}
			}
		}
		// Clear cookie by setting expired
		c.SetCookie(RefreshCookieName, "", -1, "/", "", h.Config.IsProduction(), true)
		c.JSON(http.StatusOK, gin.H{"message": "logged out"})
	}
}

func (h *AuthHandlers) HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}