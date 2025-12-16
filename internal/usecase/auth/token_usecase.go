package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	"github.com/dhanarrizky/Golang-template/internal/ports"
)

var (
	ErrRefreshTokenNotFound   = errors.New("refresh token not found")
	ErrRefreshTokenExpired    = errors.New("refresh token expired")
	ErrRefreshTokenCompromised = errors.New("possible token reuse detected")
	ErrRefreshTokenRevoked     = errors.New("refresh token has been revoked")
)

type RefreshResult struct {
	AccessToken string
	AccessExp   time.Time
	NewRefreshToken string // untuk rotation
	NewRefreshExp   time.Time
}

type TokenUsecase interface {
	Refresh(ctx context.Context, oldRefreshToken string, deviceName string) (*RefreshResult, error)
	Revoke(ctx context.Context, refreshToken string) error
	RevokeFamily(ctx context.Context, familyID string) error // jika deteksi reuse
}

type tokenUsecase struct {
	refreshRepo   ports.RefreshTokenRepository
	sessionRepo   ports.UserSessionRepository
	jwtSecret     string
	accessExp     time.Duration
	refreshExp    time.Duration
}

func NewTokenUsecase(
	refreshRepo ports.RefreshTokenRepository,
	sessionRepo ports.UserSessionRepository,
	jwtSecret string,
	accessExp, refreshExp time.Duration,
) TokenUsecase {
	return &tokenUsecase{
		refreshRepo:   refreshRepo,
		sessionRepo:   sessionRepo,
		jwtSecret:     jwtSecret,
		accessExp:     accessExp,
		refreshExp:    refreshExp,
	}
}

func (u *tokenUsecase) Refresh(ctx context.Context, oldRefreshToken string, deviceName string) (*RefreshResult, error) {
	// 1. Cari token
	token, err := u.refreshRepo.FindByToken(ctx, oldRefreshToken)
	if err != nil || token == nil {
		return nil, ErrRefreshTokenNotFound
	}

	if token.Revoked {
		// Jika sudah direvoke → kemungkinan reuse → revoke seluruh family
		u.revokeFamilyAndSessions(ctx, token.FamilyID)
		return nil, ErrRefreshTokenCompromised
	}

	if time.Now().After(token.ExpiresAt) {
		return nil, ErrRefreshTokenExpired
	}

	// 2. Cek reuse detection: apakah token ini sudah pernah dipakai untuk refresh?
	if token.Used {
		// Token reuse detected → revoke entire family
		u.revokeFamilyAndSessions(ctx, token.FamilyID)
		return nil, ErrRefreshTokenCompromised
	}

	// 3. Mark as used (untuk deteksi reuse)
	if err := u.refreshRepo.MarkAsUsed(ctx, oldRefreshToken); err != nil {
		return nil, err
	}

	// 4. Rotate: buat refresh token baru (same family)
	newRefreshToken := uuid.New().String()
	newExpiresAt := time.Now().Add(u.refreshExp)

	err = u.refreshRepo.Create(ctx, domain.RefreshToken{
		Token:     newRefreshToken,
		UserID:    token.UserID,
		FamilyID:  token.FamilyID,
		Device:    deviceName,
		ExpiresAt: newExpiresAt,
	})
	if err != nil {
		return nil, err
	}

	// 5. Update last used session
	u.sessionRepo.UpdateLastUsed(ctx, token.UserID, token.FamilyID, time.Now())

	// 6. Generate new access token
	user, err := u.refreshRepo.GetUserFromToken(ctx, oldRefreshToken) // atau inject user repo jika perlu
	if err != nil {
		return nil, err
	}

	accessExp := time.Now().Add(u.accessExp)
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"roles": user.Roles,
		"exp":   accessExp.Unix(),
		"iat":   time.Now().Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := jwtToken.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &RefreshResult{
		AccessToken:     accessToken,
		AccessExp:       accessExp,
		NewRefreshToken: newRefreshToken,
		NewRefreshExp:   newExpiresAt,
	}, nil
}

func (u *tokenUsecase) Revoke(ctx context.Context, refreshToken string) error {
	return u.refreshRepo.Revoke(ctx, refreshToken)
}

func (u *tokenUsecase) revokeFamilyAndSessions(ctx context.Context, familyID string) {
	// Revoke semua token dalam family
	u.refreshRepo.RevokeFamily(ctx, familyID)
	// Optional: hapus session terkait atau mark inactive
	u.sessionRepo.RevokeByFamily(ctx, familyID)
}