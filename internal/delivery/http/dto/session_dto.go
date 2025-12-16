package dto

import "time"

// Item session (satu device / login)
type SessionResponse struct {
	ID        string    `json:"id"`
	Device    string    `json:"device"`
	IP        string    `json:"ip"`
	LastUsed  time.Time `json:"last_used"`
	CreatedAt time.Time `json:"created_at"`
	Current   bool      `json:"current"` // apakah ini session sekarang
}

// Response list session
type ListSessionResponse struct {
	Sessions []SessionResponse `json:"sessions"`
}

// Response revoke session
type RevokeSessionResponse struct {
	Message string `json:"message"`
}
