package auth

// import (
// 	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
// 	model "github.com/dhanarrizky/Golang-template/internal/infrastructure/database/models/auth"
// )

// func ToDomainPasswordResetToken(m *model.PasswordResetToken) *domain.PasswordResetToken {
// 	if m == nil {
// 		return nil
// 	}

// 	return &domain.PasswordResetToken{
// 		ID:        m.ID,
// 		UserID:    m.UserID,
// 		TokenHash: m.TokenHash,
// 		ExpiresAt: m.ExpiresAt,
// 		Used:      m.Used,
// 		CreatedAt: m.CreatedAt,
// 	}
// }

// func ToModelPasswordResetToken(d *domain.PasswordResetToken) *model.PasswordResetToken {
// 	if d == nil {
// 		return nil
// 	}

// 	return &model.PasswordResetToken{
// 		ID:        d.ID,
// 		UserID:    d.UserID,
// 		TokenHash: d.TokenHash,
// 		ExpiresAt: d.ExpiresAt,
// 		Used:      d.Used,
// 		CreatedAt: d.CreatedAt,
// 	}
// }
