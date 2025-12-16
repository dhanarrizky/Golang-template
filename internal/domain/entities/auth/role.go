package auth

import "time"

// Role adalah domain concept untuk authorization
type Role struct {
	ID uint64

	Name        string
	Description *string

	CreatedAt time.Time
}

/* ===== Domain Behavior ===== */

func (r *Role) Rename(name string) {
	r.Name = name
}

func (r *Role) ChangeDescription(desc *string) {
	r.Description = desc
}
