package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/dhanarrizky/Golang-template/internal/delivery/http/dto"
	"github.com/dhanarrizky/Golang-template/internal/usecase/user"
)

type UserHandler struct {
	usecase  user.UserUsecase
	validate *validator.Validate
}

func NewUserHandler(usecase user.UserUsecase, validate *validator.Validate) *UserHandler {
	return &UserHandler{
		usecase:  usecase,
		validate: validate,
	}
}

// GET /users/me
func (h *UserHandler) Me(c *gin.Context) {
	userID := c.GetString("user_id")

	user, err := h.usecase.GetMe(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.UserProfileResponse{
		ID:            user.ID,
		Email:         user.Email,
		Username:      user.Username,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
	})
}

// PUT /users/me
func (h *UserHandler) Update(c *gin.Context) {
	userID := c.GetString("user_id")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Validation failed"})
		return
	}

	if err := h.usecase.UpdateProfile(c.Request.Context(), userID, req.Username); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.UpdateProfileResponse{
		Message: "Profile updated successfully",
	})
}

// DELETE /users/me
func (h *UserHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := h.usecase.SoftDelete(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, dto.DeleteAccountResponse{
		Message: "Account deleted successfully",
	})
}
