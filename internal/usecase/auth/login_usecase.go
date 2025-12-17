package auth

import (
	"context"
	"errors"
	"time"

	"github.com/dhanarrizky/Golang-template/internal/ports"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type LoginResult struct {
	AccessToken   string
	AccessExp     time.Time
	RefreshToken  string
	RefreshExp    time.Time
	UserID        string
	Email         string
	Username      string
	Roles         []string
	EmailVerified bool
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTooManyAttempts    = errors.New("too many login attempts")
	ErrAccountLocked      = errors.New("account locked")
)

type LoginUsecase interface {
	Execute(ctx context.Context, identifier, password, deviceName string) (*LoginResult, error)
}

type loginUsecase struct {
	userRepo         ports.UserRepository
	loginAttemptRepo ports.LoginAttemptRepository
	sessionRepo      ports.UserSessionRepository
	refreshRepo      ports.RefreshTokenRepository
	passwordHasher   ports.PasswordHasher
	jwtSecret        string
	accessExp        time.Duration
	refreshExp       time.Duration
}

func NewLoginUsecase(
	userRepo ports.UserRepository,
	loginAttemptRepo ports.LoginAttemptRepository,
	sessionRepo ports.UserSessionRepository,
	refreshRepo ports.RefreshTokenRepository,
	passwordHasher ports.PasswordHasher,
	jwtSecret string,
	accessExp, refreshExp time.Duration,
) LoginUsecase {
	return &loginUsecase{
		userRepo:         userRepo,
		loginAttemptRepo: loginAttemptRepo,
		sessionRepo:      sessionRepo,
		refreshRepo:      refreshRepo,
		passwordHasher:   passwordHasher,
		jwtSecret:        jwtSecret,
		accessExp:        accessExp,
		refreshExp:       refreshExp,
	}
}

func (u *loginUsecase) Execute(ctx context.Context, identifier, password, deviceName string) (*LoginResult, error) {
	// 1. Check rate limit
	if u.loginAttemptRepo.IsRateLimited(ctx, identifier) {
		return nil, ErrTooManyAttempts
	}

	// 2. Find user
	user, err := u.userRepo.FindByEmailOrUsername(ctx, identifier)
	if err != nil || user == nil || !u.passwordHasher.Compare(password, user.HashedPassword) {
		u.loginAttemptRepo.RecordFailedAttempt(ctx, identifier)
		return nil, ErrInvalidCredentials
	}

	if user.Locked {
		return nil, ErrAccountLocked
	}

	// 3. Reset attempt on success
	u.loginAttemptRepo.ResetAttempts(ctx, identifier)

	// 4. Generate refresh token & family
	familyID := uuid.New().String()
	refreshToken := uuid.New().String()

	err = u.refreshRepo.Create(ctx, domain.RefreshToken{
		Token:     refreshToken,
		UserID:    user.ID,
		FamilyID:  familyID,
		Device:    deviceName,
		ExpiresAt: time.Now().Add(u.refreshExp),
	})
	if err != nil {
		return nil, err
	}

	// 5. Create session
	u.sessionRepo.Create(ctx, domain.UserSession{
		UserID:   user.ID,
		Device:   deviceName,
		FamilyID: familyID,
		LastUsed: time.Now(),
	})

	// 6. Generate JWT access token
	accessExp := time.Now().Add(u.accessExp)
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"roles": user.Roles,
		"exp":   accessExp.Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:   accessToken,
		AccessExp:     accessExp,
		RefreshToken:  refreshToken,
		RefreshExp:    time.Now().Add(u.refreshExp),
		UserID:        user.ID,
		Email:         user.Email,
		Username:      user.Username,
		Roles:         user.Roles,
		EmailVerified: user.EmailVerified,
	}, nil
}
