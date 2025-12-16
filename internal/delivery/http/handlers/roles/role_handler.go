package roles

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/dhanarrizky/Golang-template/internal/delivery/http/dto"
	"github.com/dhanarrizky/Golang-template/internal/usecase/role"
)

type RoleHandler struct {
	usecase  role.RoleUsecase
	validate *validator.Validate
}

func NewRoleHandler(usecase role.RoleUsecase, validate *validator.Validate) *RoleHandler {
	return &RoleHandler{
		usecase:  usecase,
		validate: validate,
	}
}

// GET /roles
func (h *RoleHandler) List(c *gin.Context) {
	roles, err := h.usecase.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "Failed to fetch roles"})
		return
	}

	resp := make([]dto.RoleResponse, 0, len(roles))
	for _, r := range roles {
		resp = append(resp, dto.RoleResponse{
			ID:        r.ID,
			Name:      r.Name,
			CreatedAt: r.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, dto.ListRoleResponse{Roles: resp})
}

// POST /roles
func (h *RoleHandler) Create(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Validation failed"})
		return
	}

	if err := h.usecase.Create(c.Request.Context(), req.Name); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.MessageResponse{Message: "Role created"})
}

// PUT /roles/:id
func (h *RoleHandler) Update(c *gin.Context) {
	roleID := c.Param("id")

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Validation failed"})
		return
	}

	if err := h.usecase.Update(c.Request.Context(), roleID, req.Name); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Role updated"})
}

// POST /roles/assign
func (h *RoleHandler) Assign(c *gin.Context) {
	var req dto.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Invalid request"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "Validation failed"})
		return
	}

	if err := h.usecase.AssignToUser(c.Request.Context(), req.UserID, req.RoleID); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Role assigned to user"})
}
