package dto

import "time"

// ===== RESPONSE =====

type RoleResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ListRoleResponse struct {
	Roles []RoleResponse `json:"roles"`
}

// ===== REQUEST =====

type CreateRoleRequest struct {
	Name string `json:"name" validate:"required,min=3,max=50"`
}

type UpdateRoleRequest struct {
	Name string `json:"name" validate:"required,min=3,max=50"`
}

type AssignRoleRequest struct {
	UserID string `json:"user_id" validate:"required,uuid4"`
	RoleID string `json:"role_id" validate:"required,uuid4"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
