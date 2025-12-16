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
