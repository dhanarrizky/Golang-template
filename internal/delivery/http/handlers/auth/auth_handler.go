package auth

import (
	"net/http"

	"github.com/dhanarrizky/Golang-template/internal/delivery/http/dto"
	"github.com/dhanarrizky/Golang-template/internal/usecase/auth"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	loginUsecase auth.LoginUsecase
	validate     *validator.Validate
}

func NewAuthHandler(
	loginUsecase auth.LoginUsecase,
	validate *validator.Validate,
) *AuthHandler {
	return &AuthHandler{
		loginUsecase: loginUsecase,
		validate:     validate,
	}
}

// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
		})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		errors := map[string]string{}
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = e.Tag()
		}

		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Validation failed",
			Errors:  errors,
		})
		return
	}

	result, err := h.loginUsecase.Login(
		c.Request.Context(),
		req.Identifier,
		req.Password,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	// Refresh token = HTTP concern → BOLEH di handler
	c.SetCookie(
		"refresh_token",
		result.RefreshToken,
		result.RefreshExp.Second(),
		"/auth",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken: result.AccessToken,
		ExpiresAt:   result.AccessExp,
		User: dto.UserInfo{
			ID:            result.UserID,
			Email:         result.Email,
			Username:      result.Username,
			Roles:         result.Roles,
			EmailVerified: result.EmailVerified,
		},
	})
}

// POST /auth/logout (current device)
func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "No refresh token",
		})
		return
	}

	// h.logoutUsecase.RevokeCurrent(...)
	c.SetCookie("refresh_token", "", -1, "/auth", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out",
	})
}

// POST /auth/logout-all
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	// extract user dari access token (middleware)
	// userID := c.GetString("user_id") // diasumsikan middleware set ini
	// panggil usecase revoke all sessions untuk userID
	c.JSON(http.StatusOK, gin.H{"message": "All sessions revoked"})
}

// GET /auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetString("user_id")

	c.JSON(http.StatusOK, dto.MeResponse{
		User: dto.UserInfo{
			ID: userID,
		},
	})
}
