package user

import (
	"context"
	"errors"

	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	authPorts "github.com/dhanarrizky/Golang-template/internal/ports/auth"
	mailerPorts "github.com/dhanarrizky/Golang-template/internal/ports/mailer"
	userPorts "github.com/dhanarrizky/Golang-template/internal/ports/users"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUsernameTaken   = errors.New("username already taken")
	ErrEmailTaken      = errors.New("email already taken")
	ErrInvalidPassword = errors.New("invalid password")
)

type UserUsecase interface {
	Register(ctx context.Context, username, email, password string) (*domain.User, error)
	GetMe(ctx context.Context, userID string) (*domain.User, error)
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID, username string) error
	UpdateUser(ctx context.Context, userID, username, email string) error // admin/full update
	SoftDelete(ctx context.Context, userID string) error
	PermanentDelete(ctx context.Context, userID string) error // admin
}

type userUsecase struct {
	userRepo    userPorts.UserRepository
	mailRepo    mailerPorts.EmailVerificationTokenRepository
	sessionRepo userPorts.UserSessionRepository
	refreshRepo authPorts.RefreshTokenRepository
}

func NewUserUsecase(
	userRepo userPorts.UserRepository,
	sessionRepo userPorts.UserSessionRepository,
	refreshRepo authPorts.RefreshTokenRepository,
) UserUsecase {
	return &userUsecase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		refreshRepo: refreshRepo,
	}
}

// ================= REGISTER =================
func (u *userUsecase) Register(ctx context.Context, username, email, password string) (*domain.User, error) {
	if username == "" || email == "" || password == "" {
		return nil, errors.New("username, email, and password are required")
	}

	// Check uniqueness
	if exists, _ := u.userRepo.ExistsByUsername(ctx, username); exists {
		return nil, ErrUsernameTaken
	}
	if exists, _ := u.userRepo.ExistsByEmail(ctx, email); exists {
		return nil, ErrEmailTaken
	}

	// Hash password using bcrypt (bukan dari PasswordHasher port)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword), // asumsikan field di domain adalah Password atau HashedPassword
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Kosongkan password sebelum return
	user.Password = ""
	return user, nil
}

// ================= GET ME =================
func (u *userUsecase) GetMe(ctx context.Context, userID string) (*domain.User, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}
	user.Password = ""
	return user, nil
}

// ================= GET BY ID =================
func (u *userUsecase) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}
	user.Password = ""
	return user, nil
}

// ================= UPDATE PROFILE (self) =================
func (u *userUsecase) UpdateProfile(ctx context.Context, userID, username string) error {
	if username == "" {
		return nil
	}

	exists, err := u.userRepo.ExistsByUsernameExceptID(ctx, username, userID)
	if err != nil || exists {
		return ErrUsernameTaken
	}

	return u.userRepo.UpdateUsername(ctx, userID, username)
}

// ================= UPDATE USER (admin/full) =================
func (u *userUsecase) UpdateUser(ctx context.Context, userID, username, email string) error {
	if username == "" && email == "" {
		return nil
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return ErrUserNotFound
	}

	if username != "" && username != user.Username {
		if exists, _ := u.userRepo.ExistsByUsernameExceptID(ctx, username, userID); exists {
			return ErrUsernameTaken
		}
		user.Username = username
	}

	if email != "" && email != user.Email {
		if exists, _ := u.userRepo.ExistsByEmailExceptID(ctx, email, userID); exists {
			return ErrEmailTaken
		}
		user.Email = email
	}

	return u.userRepo.Update(ctx, user)
}

// ================= SOFT DELETE =================
func (u *userUsecase) SoftDelete(ctx context.Context, userID string) error {
	if err := u.userRepo.SoftDelete(ctx, userID); err != nil {
		return err
	}

	// Revoke all sessions & tokens
	u.refreshRepo.RevokeAllByUser(ctx, userID)
	u.sessionRepo.RevokeAllByUser(ctx, userID)

	return nil
}

// ================= PERMANENT DELETE (admin) =================
func (u *userUsecase) PermanentDelete(ctx context.Context, userID string) error {
	u.refreshRepo.RevokeAllByUser(ctx, userID)
	u.sessionRepo.RevokeAllByUser(ctx, userID)

	return u.userRepo.PermanentDelete(ctx, userID)
}
