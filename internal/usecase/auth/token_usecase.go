package auth

import (
	"context"
	"errors"
	"time"

	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	"github.com/dhanarrizky/Golang-template/internal/ports"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrRefreshTokenNotFound    = errors.New("refresh token not found")
	ErrRefreshTokenExpired     = errors.New("refresh token expired")
	ErrRefreshTokenCompromised = errors.New("possible token reuse detected")
)

type LoginTokenResult struct {
	AccessToken  string
	AccessExp    time.Time
	RefreshToken string
	RefreshExp   time.Time
}

type RefreshResult struct {
	AccessToken     string
	AccessExp       time.Time
	NewRefreshToken string
	NewRefreshExp   time.Time
}

type TokenUsecase interface {
	IssueForLogin(
		ctx context.Context,
		user domain.User,
		deviceName string,
	) (*LoginTokenResult, error)

	Refresh(
		ctx context.Context,
		oldRefreshToken string,
		deviceName string,
	) (*RefreshResult, error)

	Revoke(ctx context.Context, refreshToken string) error
	RevokeFamily(ctx context.Context, familyID string) error
}

type tokenUsecase struct {
	refreshRepo ports.RefreshTokenRepository
	sessionRepo ports.UserSessionRepository
	userRepo    ports.UserRepository
	jwtSecret   string
	accessExp   time.Duration
	refreshExp  time.Duration
}

func NewTokenUsecase(
	refreshRepo ports.RefreshTokenRepository,
	sessionRepo ports.UserSessionRepository,
	userRepo ports.UserRepository,
	jwtSecret string,
	accessExp, refreshExp time.Duration,
) TokenUsecase {
	return &tokenUsecase{
		refreshRepo: refreshRepo,
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		accessExp:   accessExp,
		refreshExp:  refreshExp,
	}
}

// ================= ISSUE TOKEN (LOGIN) =================

func (u *tokenUsecase) IssueForLogin(
	ctx context.Context,
	user domain.User,
	deviceName string,
) (*LoginTokenResult, error) {

	familyID := uuid.New().String()
	refreshToken := uuid.New().String()
	refreshExp := time.Now().Add(u.refreshExp)

	err := u.refreshRepo.Create(ctx, domain.RefreshToken{
		TokenHash: refreshToken,
		UserID:    user.ID,
		FamilyID:  familyID,
		IPAddress: deviceName,
		ExpiresAt: refreshExp,
	})
	if err != nil {
		return nil, err
	}

	u.sessionRepo.Create(ctx, domain.UserSession{
		UserID:    user.ID,
		FamilyID:  familyID,
		IPAddress: deviceName,
		LastUsed:  time.Now(),
	})

	accessToken, accessExp, err := u.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	return &LoginTokenResult{
		AccessToken:  accessToken,
		AccessExp:    accessExp,
		RefreshToken: refreshToken,
		RefreshExp:   refreshExp,
	}, nil
}

// ================= REFRESH TOKEN =================

func (u *tokenUsecase) Refresh(
	ctx context.Context,
	oldRefreshToken string,
	deviceName string,
) (*RefreshResult, error) {

	token, err := u.refreshRepo.FindByToken(ctx, oldRefreshToken)
	if err != nil || token == nil {
		return nil, ErrRefreshTokenNotFound
	}

	if token.Revoked || token.Used {
		u.RevokeFamily(ctx, token.FamilyID)
		return nil, ErrRefreshTokenCompromised
	}

	if time.Now().After(token.ExpiresAt) {
		return nil, ErrRefreshTokenExpired
	}

	if err := u.refreshRepo.MarkAsUsed(ctx, oldRefreshToken); err != nil {
		return nil, err
	}

	newRefreshToken := uuid.New().String()
	newRefreshExp := time.Now().Add(u.refreshExp)

	err = u.refreshRepo.Create(ctx, domain.RefreshToken{
		TokenHash: newRefreshToken,
		UserID:    token.UserID,
		FamilyID:  token.FamilyID,
		UserAgent: deviceName,
		IPAddress: deviceName,
		ExpiresAt: newRefreshExp,
	})
	if err != nil {
		return nil, err
	}

	u.sessionRepo.UpdateLastUsed(ctx, token.UserID, token.FamilyID, time.Now())

	user, err := u.userRepo.FindByID(ctx, token.UserID)
	if err != nil {
		return nil, err
	}

	accessToken, accessExp, err := u.generateAccessToken(*user)
	if err != nil {
		return nil, err
	}

	return &RefreshResult{
		AccessToken:     accessToken,
		AccessExp:       accessExp,
		NewRefreshToken: newRefreshToken,
		NewRefreshExp:   newRefreshExp,
	}, nil
}

// ================= REVOKE =================

func (u *tokenUsecase) Revoke(ctx context.Context, refreshToken string) error {
	return u.refreshRepo.Revoke(ctx, refreshToken)
}

func (u *tokenUsecase) RevokeFamily(ctx context.Context, familyID string) error {
	u.refreshRepo.RevokeFamily(ctx, familyID)
	u.sessionRepo.RevokeByFamily(ctx, familyID)
	return nil
}

// ================= HELPERS =================

func (u *tokenUsecase) generateAccessToken(
	user domain.User,
) (string, time.Time, error) {

	exp := time.Now().Add(u.accessExp)

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"roles": user.RoleID,
		"exp":   exp.Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, exp, nil
}
