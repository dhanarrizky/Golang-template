package auth

import (
	"context"
	"errors"
	"time"

	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	authPort "github.com/dhanarrizky/Golang-template/internal/ports/auth"
	userPort "github.com/dhanarrizky/Golang-template/internal/ports/users"
	"github.com/dhanarrizky/Golang-template/pkg/utils"
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
	refreshRepo           authPort.RefreshTokenRepository
	sessionRepo           userPort.UserSessionRepository
	userRepo              userPort.UserRepository
	accessExp, refreshExp time.Duration
	tokenSigner           authPort.TokenSigner
}

func NewTokenUsecase(
	refreshRepo authPort.RefreshTokenRepository,
	sessionRepo userPort.UserSessionRepository,
	userRepo userPort.UserRepository,
	accessExp, refreshExp time.Duration,
	tokenSigner authPort.TokenSigner,
) TokenUsecase {
	return &tokenUsecase{
		refreshRepo: refreshRepo,
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
		accessExp:   accessExp,
		refreshExp:  refreshExp,
		tokenSigner: tokenSigner,
	}
}

// ================= ISSUE TOKEN (LOGIN) =================
func (u *tokenUsecase) IssueForLogin(
	ctx context.Context,
	user domain.User,
	deviceName string,
) (*LoginTokenResult, error) {

	familyID, err := utils.GenerateID()
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.tokenSigner.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshExp := time.Now().Add(u.refreshExp)

	err = u.refreshRepo.Create(ctx, domain.RefreshToken{
		TokenHash: refreshToken,
		UserID:    user.ID,
		FamilyID:  familyID,
		IPAddress: deviceName,
		ExpiresAt: refreshExp,
	})
	if err != nil {
		return nil, err
	}

	now := time.Now()
	_ = u.sessionRepo.Create(ctx, domain.UserSession{
		UserID:     user.ID,
		FamilyID:   familyID,
		IPAddress:  deviceName,
		LastSeenAt: &now,
	})

	claims := map[string]any{
		"email": user.Email,
		"roles": user.RoleID,
	}

	accessToken, err := u.tokenSigner.GenerateAccessToken(user.ID, claims)
	if err != nil {
		return nil, err
	}

	return &LoginTokenResult{
		AccessToken:  accessToken,
		AccessExp:    time.Now().Add(u.accessExp),
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

	token, err := u.refreshRepo.GetByTokenHash(ctx, oldRefreshToken)
	if err != nil || token == nil {
		return nil, ErrRefreshTokenNotFound
	}

	if token.Revoked || token.Used {
		_ = u.RevokeFamily(ctx, token.FamilyID)
		return nil, ErrRefreshTokenCompromised
	}

	if time.Now().After(token.ExpiresAt) {
		return nil, ErrRefreshTokenExpired
	}

	if err := u.refreshRepo.MarkAsUsed(ctx, oldRefreshToken); err != nil {
		return nil, err
	}

	newRefreshToken, err := u.tokenSigner.GenerateRefreshToken(token.UserID)
	if err != nil {
		return nil, err
	}

	newRefreshExp := time.Now().Add(u.refreshExp)

	err = u.refreshRepo.Create(ctx, domain.RefreshToken{
		TokenHash: newRefreshToken,
		UserID:    token.UserID,
		FamilyID:  token.FamilyID,
		IPAddress: deviceName,
		ExpiresAt: newRefreshExp,
	})
	if err != nil {
		return nil, err
	}

	_ = u.sessionRepo.UpdateLastUsed(ctx, token.UserID, token.FamilyID, time.Now())

	user, err := u.userRepo.FindByID(ctx, token.UserID)
	if err != nil {
		return nil, err
	}

	claims := map[string]any{
		"email": user.Email,
		"roles": user.RoleID,
	}

	accessToken, err := u.tokenSigner.GenerateAccessToken(user.ID, claims)
	if err != nil {
		return nil, err
	}

	return &RefreshResult{
		AccessToken:     accessToken,
		AccessExp:       time.Now().Add(u.accessExp),
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
