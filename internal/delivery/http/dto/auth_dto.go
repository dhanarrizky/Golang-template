package dto

import "time"

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required,min=3,max=100"` // email or username
	Password   string `json:"password" validate:"required,min=8,max=100"`
	DeviceName string `json:"device_name,omitempty" validate:"max=100"` // optional
}

type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	User        UserInfo  `json:"user"`
}

type UserInfo struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	Username  string   `json:"username,omitempty"`
	Roles     []string `json:"roles,omitempty"`
	EmailVerified bool  `json:"email_verified"`
}

type MeResponse struct {
	User UserInfo `json:"user"`
}

type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"` // untuk validation detail
}

// Tambahkan ini di bawah DTO yang sudah ada

type RefreshRequest struct {
	// Kosong karena refresh token diambil dari HttpOnly cookie
}

type RefreshResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type RevokeRequest struct {
	// Opsional: jika ingin revoke token lain, bukan current
}

type RevokeResponse struct {
	Message string `json:"message"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required,min=32"` // panjang token cukup besar
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=100"`
}

type ChangePasswordResponse struct {
	Message string `json:"message"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required,min=32"`
}

type VerifyEmailResponse struct {
	Message string `json:"message"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResendVerificationResponse struct {
	Message string `json:"message"`
}
