package auth

import (
	"net/http"
	"time"

	"github.com/dhanarrizky/Golang-template/internal/delivery/http/dto"
	"github.com/dhanarrizky/Golang-template/internal/usecase/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	loginUsecase auth.LoginUsecase
	validate     *validator.Validate
	jwtSecret    string // atau inject JWT service jika lebih kompleks
	accessExp    time.Duration
}

func NewAuthHandler(loginUsecase auth.LoginUsecase, validate *validator.Validate, jwtSecret string, accessExp time.Duration) *AuthHandler {
	return &AuthHandler{
		loginUsecase: loginUsecase,
		validate:     validate,
		jwtSecret:    jwtSecret,
		accessExp:    accessExp,
	}
}

// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = e.Tag()
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Validation failed", Errors: errors})
		return
	}

	result, err := h.loginUsecase.Execute(c.Request.Context(), req.Identifier, req.Password, req.DeviceName)
	if err != nil {
		// handle different error types (misalnya auth.ErrInvalidCredentials, auth.ErrTooManyAttempts, dll)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: err.Error()})
		return
	}

	// Set HttpOnly cookie untuk refresh token
	c.SetCookie(
		"refresh_token",
		result.RefreshToken,
		int(result.RefreshExp.Sub(time.Now()).Seconds()),
		"/auth",
		"", // domain kosong = current
		true,  // Secure (true di production)
		true,  // HttpOnly
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
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "No refresh token"})
		return
	}

	// panggil usecase logout current (revoke refresh token ini)
	// misalnya h.logoutUsecase.RevokeCurrent(...)
	c.SetCookie("refresh_token", "", -1, "/auth", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

// POST /auth/logout-all
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	// extract user dari access token (middleware)
	userID := c.GetString("user_id") // diasumsikan middleware set ini
	// panggil usecase revoke all sessions untuk userID
	c.JSON(http.StatusOK, gin.H{"message": "All sessions revoked"})
}

// GET /auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetString("user_id")
	// ambil profile dari usecase atau langsung dari claims jika cukup
	c.JSON(http.StatusOK, dto.MeResponse{
		User: dto.UserInfo{
			ID: userID,
			// isi lain dari context atau query repo jika perlu
		},
	})
}

// Helper untuk extract user dari access token (dipakai middleware)
func (h *AuthHandler) ExtractClaims(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}