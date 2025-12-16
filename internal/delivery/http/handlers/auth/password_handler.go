package auth

import (
	"net/http"

	"github.com/dhanarrizky/Golang-template/internal/delivery/http/dto"
	"github.com/dhanarrizky/Golang-template/internal/usecase/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PasswordHandler struct {
	passwordUsecase auth.PasswordUsecase
	validate        *validator.Validate
}

func NewPasswordHandler(passwordUsecase auth.PasswordUsecase, validate *validator.Validate) *PasswordHandler {
	return &PasswordHandler{
		passwordUsecase: passwordUsecase,
		validate:        validate,
	}
}

// POST /auth/password/forgot
func (h *PasswordHandler) Forgot(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Validation failed"})
		return
	}

	// Selalu return success (anti-enumeration)
	if err := h.passwordUsecase.Forgot(c.Request.Context(), req.Email); err != nil {
		// Log error internally
	}

	c.JSON(http.StatusOK, dto.ForgotPasswordResponse{Message: "If the email exists, a reset link has been sent"})
}

// POST /auth/password/reset
func (h *PasswordHandler) Reset(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Validation failed"})
		return
	}

	err := h.passwordUsecase.Reset(c.Request.Context(), req.Token, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ResetPasswordResponse{Message: "Password successfully reset. Please login again."})
}

// POST /auth/password/change (require auth middleware)
func (h *PasswordHandler) Change(c *gin.Context) {
	userID := c.GetString("user_id") // dari middleware auth

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Validation failed"})
		return
	}

	err := h.passwordUsecase.Change(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ChangePasswordResponse{Message: "Password changed successfully"})
}