package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/dhanarrizky/Golang-template/internal/delivery/http/dto"
	"github.com/dhanarrizky/Golang-template/internal/usecase/auth"
)

type EmailHandler struct {
	usecase  auth.EmailUsecase
	validate *validator.Validate
}

func NewEmailHandler(usecase auth.EmailUsecase, validate *validator.Validate) *EmailHandler {
	return &EmailHandler{
		usecase:  usecase,
		validate: validate,
	}
}

// POST /auth/email/verify
func (h *EmailHandler) Verify(c *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Validation failed"})
		return
	}

	if err := h.usecase.Verify(c.Request.Context(), req.Token); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.VerifyEmailResponse{
		Message: "Email successfully verified",
	})
}

// POST /auth/email/resend
func (h *EmailHandler) Resend(c *gin.Context) {
	var req dto.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Validation failed"})
		return
	}

	// Silent (anti enumeration)
	h.usecase.Resend(c.Request.Context(), req.Email)

	c.JSON(http.StatusOK, dto.ResendVerificationResponse{
		Message: "If the email exists, a verification link has been sent",
	})
}
