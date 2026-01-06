package user

import (
	"context"
	"errors"

	"github.com/dhanarrizky/Golang-template/internal/delivery/http/dto"
	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	authPorts "github.com/dhanarrizky/Golang-template/internal/ports/auth"
	otherPorts "github.com/dhanarrizky/Golang-template/internal/ports/others"
	userPorts "github.com/dhanarrizky/Golang-template/internal/ports/users"
)

var (
	ErrDecode          = errors.New("internal server decode")
	ErrUserNotFound    = errors.New("user not found")
	ErrUsernameTaken   = errors.New("username already taken")
	ErrEmailTaken      = errors.New("email already taken")
	ErrInvalidPassword = errors.New("invalid password")
)

// type UserUsecase interface {
// 	Register(ctx context.Context, username, email, password string) (*dto.CreateUserResponse, error)
// 	GetMe(ctx context.Context, userID string) (*dto.UserProfileResponse, error)
// 	GetUserByID(ctx context.Context, userID string) (*dto.UserResponse, error)
// 	UpdateProfile(ctx context.Context, userID, username string) error
// 	UpdateUser(ctx context.Context, userID, username, email string) error // admin/full update
// 	SoftDelete(ctx context.Context, userID string) error
// 	PermanentDelete(ctx context.Context, userID string) error // admin
// }

type UserUsecase interface {
	Register(ctx context.Context, username, email, password string) (*dto.CreateUserResponse, error)
	GetMe(ctx context.Context, userID string) (*dto.UserProfileResponse, error)
	GetUserByID(ctx context.Context, userID string) (*dto.UserResponse, error)

	// 🔹 NEW
	GetList(ctx context.Context) ([]dto.UserResponse, error)

	UpdateProfile(ctx context.Context, userID, username string) error
	UpdateUser(ctx context.Context, userID, username, email string) error
	SoftDelete(ctx context.Context, userID string) error
	PermanentDelete(ctx context.Context, userID string) error
}

type userUsecase struct {
	userRepo          userPorts.UserRepository
	sessionRepo       userPorts.UserSessionRepository
	passwordHasher    userPorts.PasswordHasher
	idCodec           otherPorts.PublicIDCodec
	refreshRepo       authPorts.RefreshTokenRepository
	refreshFamilyRepo authPorts.RefreshTokenFamilyRepository
}

func NewUserUsecase(
	userRepo userPorts.UserRepository,
	sessionRepo userPorts.UserSessionRepository,
	passwordHasher userPorts.PasswordHasher,
	idCodec otherPorts.PublicIDCodec,
	refreshRepo authPorts.RefreshTokenRepository,
	refreshFamilyRepo authPorts.RefreshTokenFamilyRepository,
) UserUsecase {
	return &userUsecase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		refreshRepo: refreshRepo,
	}
}

// ================= REGISTER =================
func (u *userUsecase) Register(ctx context.Context, username, email, password string) (*dto.CreateUserResponse, error) {
	if username == "" || email == "" || password == "" {
		return nil, errors.New("username, email, and password are required")
	}

	// Check uniqueness
	if exists, _ := u.userRepo.GetByEmailOrUsername(ctx, username); exists != nil {
		return nil, ErrUsernameTaken
	}
	if exists, _ := u.userRepo.GetByEmail(ctx, email); exists != nil {
		return nil, ErrEmailTaken
	}

	// Hash password using bcrypt (bukan dari PasswordHasher port)
	hashedPassword, err := u.passwordHasher.HashPassword([]byte(password))
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword, // asumsikan field di domain adalah Password atau HashedPassword
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Kosongkan password sebelum return
	encrypId, err := u.idCodec.Encode(user.ID)
	if err != nil {
		return nil, err
	}

	result := dto.CreateUserResponse{
		ID:       encrypId,
		Email:    user.Email,
		Username: user.Username,
	}

	return &result, nil
}

// ================= GET ME =================
func (u *userUsecase) GetMe(ctx context.Context, userID string) (*dto.UserProfileResponse, error) {
	id, err := u.idCodec.Decode(userID)
	if err != nil {
		return nil, ErrDecode
	}

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	encrypId, err := u.idCodec.Encode(user.ID)
	if err != nil {
		return nil, err
	}

	result := dto.UserProfileResponse{
		ID:            encrypId,
		Email:         user.Email,
		Username:      user.Username,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
	}

	return &result, nil
}

// ================= GET BY ID =================
func (u *userUsecase) GetUserByID(ctx context.Context, userID string) (*dto.UserResponse, error) {
	id, err := u.idCodec.Decode(userID)
	if err != nil {
		return nil, ErrDecode
	}

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	encrypId, err := u.idCodec.Encode(user.ID)
	if err != nil {
		return nil, err
	}

	result := dto.UserResponse{
		ID:        encrypId,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}

	return &result, nil
}

// ================= GET LIST (admin) =================
func (u *userUsecase) GetList(ctx context.Context) ([]dto.UserResponse, error) {

	users, err := u.userRepo.GetList(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.UserResponse, 0, len(users))
	for _, usr := range users {
		encrypId, err := u.idCodec.Encode(usr.ID)
		if err != nil {
			return nil, err
		}

		result = append(result, dto.UserResponse{
			ID:        encrypId,
			Email:     usr.Email,
			Username:  usr.Username,
			CreatedAt: usr.CreatedAt,
		})
	}

	return result, nil
}

// ================= UPDATE PROFILE (self) =================
func (u *userUsecase) UpdateProfile(ctx context.Context, userID, username string) error {
	if username == "" {
		return nil
	}

	id, err := u.idCodec.Decode(userID)
	if err != nil {
		return ErrDecode
	}

	exists, err := u.userRepo.ExistsByUsernameExceptID(ctx, username, id)
	if err != nil || exists {
		return ErrUsernameTaken
	}

	return u.userRepo.UpdateUsername(ctx, id, username)
}

// ================= UPDATE USER (admin/full) =================
func (u *userUsecase) UpdateUser(ctx context.Context, userID, username, email string) error {
	if username == "" && email == "" {
		return nil
	}

	id, err := u.idCodec.Decode(userID)
	if err != nil {
		return ErrDecode
	}

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil || user == nil {
		return ErrUserNotFound
	}

	if username != "" && username != user.Username {
		if exists, _ := u.userRepo.ExistsByUsernameExceptID(ctx, username, id); exists {
			return ErrUsernameTaken
		}
		user.Username = username
	}

	if email != "" && email != user.Email {
		if exists, _ := u.userRepo.ExistsByEmailExceptID(ctx, email, id); exists {
			return ErrEmailTaken
		}
		user.Email = email
	}

	return u.userRepo.Update(ctx, user)
}

// ================= SOFT DELETE =================
func (u *userUsecase) SoftDelete(ctx context.Context, userID string) error {
	id, err := u.idCodec.Decode(userID)
	if err != nil {
		return ErrDecode
	}

	families, err := u.refreshFamilyRepo.GetByUserID(ctx, id)
	if err != nil {
		return err
	}

	for _, family := range families {
		_ = u.refreshRepo.RevokeByFamily(ctx, family.ID)
		_ = u.refreshFamilyRepo.Revoke(ctx, family.ID)
	}

	err = u.userRepo.SoftDelete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// ================= PERMANENT DELETE (admin) =================
func (u *userUsecase) PermanentDelete(ctx context.Context, userID string) error {
	id, err := u.idCodec.Decode(userID)
	if err != nil {
		return ErrDecode
	}

	families, err := u.refreshFamilyRepo.GetByUserID(ctx, id)
	if err != nil {
		return err
	}

	for _, family := range families {
		_ = u.refreshRepo.RevokeByFamily(ctx, family.ID)
		_ = u.refreshFamilyRepo.Revoke(ctx, family.ID)
	}

	err = u.userRepo.SoftDelete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
