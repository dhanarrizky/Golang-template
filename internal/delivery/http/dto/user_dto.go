package dto

import "time"

type UserProfileResponse struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
}

type UpdateProfileRequest struct {
	Username string `json:"username" validate:"omitempty,min=3,max=50"`
}

type UpdateProfileResponse struct {
	Message string `json:"message"`
}

type DeleteAccountResponse struct {
	Message string `json:"message"`
}

// Create
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=8"`
}

type CreateUserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// GetByID
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// List
type UserListResponse struct {
	Items []UserResponse `json:"items"`
}

// Delete
// type MessageResponse struct {
// 	Message string `json:"message"`
// }
