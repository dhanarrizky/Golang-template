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

// GET /users
func (h *UserHandler) List(c *gin.Context) {
	users, err := h.usecase.get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to fetch users",
		})
		return
	}

	c.JSON(http.StatusOK, dto.UserListResponse{
		Items: users,
	})
}

// POST /users
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Validation failed"})
		return
	}

	user, err := h.usecase.Register(
		c.Request.Context(),
		req.Email,
		req.Username,
		req.Password,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GET /users/me
func (h *UserHandler) Me(c *gin.Context) {
	userID := c.GetString("user_id")

	user, err := h.usecase.GetMe(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GET /users/:id
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.usecase.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// POST /users/me/verify-email
func (h *UserHandler) VerifyEmail(c *gin.Context) {
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

// DELETE /users/:id/permanent
func (h *UserHandler) PermanentDelete(c *gin.Context) {
	id := c.Param("id")

	if err := h.usecase.PermanentDelete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to permanently delete user",
		})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "User permanently deleted",
	})
}
