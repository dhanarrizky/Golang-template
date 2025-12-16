package auth

import (
	"net/http"
	"time"

	"github.com/dhanarrizky/Golang-template/internal/delivery/http/dto"
	"github.com/dhanarrizky/Golang-template/internal/usecase/auth"

	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	tokenUsecase auth.TokenUsecase
	jwtSecret    string
}

func NewTokenHandler(tokenUsecase auth.TokenUsecase, jwtSecret string) *TokenHandler {
	return &TokenHandler{
		tokenUsecase: tokenUsecase,
		jwtSecret:    jwtSecret,
	}
}

// POST /auth/refresh
func (h *TokenHandler) Refresh(c *gin.Context) {
	oldRefreshToken, err := c.Cookie("refresh_token")
	if err != nil || oldRefreshToken == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: "Missing refresh token"})
		return
	}

	// Optional: ambil device name dari header atau context
	deviceName := c.GetHeader("User-Agent") // atau dari session sebelumnya

	result, err := h.tokenUsecase.Refresh(c.Request.Context(), oldRefreshToken, deviceName)
	if err != nil {
		// Jika compromised, clear cookie
		if err == auth.ErrRefreshTokenCompromised {
			c.SetCookie("refresh_token", "", -1, "/auth", "", true, true)
		}
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: err.Error()})
		return
	}

	// Set new rotated refresh token (HttpOnly cookie)
	c.SetCookie(
		"refresh_token",
		result.NewRefreshToken,
		int(result.NewRefreshExp.Sub(time.Now()).Seconds()),
		"/auth",
		"",
		true,  // Secure
		true,  // HttpOnly
	)

	c.JSON(http.StatusOK, dto.RefreshResponse{
		AccessToken: result.AccessToken,
		ExpiresAt:   result.AccessExp,
	})
}

// POST /auth/revoke (opsional - revoke current refresh token)
func (h *TokenHandler) Revoke(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "No refresh token"})
		return
	}

	if err := h.tokenUsecase.Revoke(c.Request.Context(), refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to revoke"})
		return
	}

	c.SetCookie("refresh_token", "", -1, "/auth", "", true, true)
	c.JSON(http.StatusOK, dto.RevokeResponse{Message: "Refresh token revoked"})
}