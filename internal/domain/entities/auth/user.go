package auth

import "time"

type User struct {
	ID uint64

	Username      string
	Email         string
	EmailVerified bool
	PasswordHash  string
	Name          *string

	RoleID uint64
	Locked bool

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

/* ===== Domain Behavior ===== */

func (u *User) VerifyEmail() {
	u.EmailVerified = true
}

func (u *User) ChangePassword(hash string) {
	u.PasswordHash = hash
}

func (u *User) SoftDelete(now time.Time) {
	u.DeletedAt = &now
}

func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}
