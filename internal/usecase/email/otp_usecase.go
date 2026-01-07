package auth

import (
	"context"
	"time"

	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
	ports "github.com/dhanarrizky/Golang-template/internal/ports/email"
)

type OTPUsecase struct {
	emailService ports.EmailSender
	otpRepo      ports.EmailOTPRepository
	expiryMinute int
}

func NewOTPUsecase(
	email ports.EmailSender,
	repo ports.EmailOTPRepository,
	expiry int,
) *OTPUsecase {
	return &OTPUsecase{
		emailService: email,
		otpRepo:      repo,
		expiryMinute: expiry,
	}
}

func (u *OTPUsecase) RequestOTP(
	ctx context.Context,
	email string,
	otp string,
	hash string,
	ip string,
	ua string,
) error {

	// 1️⃣ Kirim email
	if err := u.emailService.SendOTP(email, otp); err != nil {
		return err
	}

	// 2️⃣ Simpan OTP sebagai ENTITY
	entity := &domain.EmailOTP{
		Email:   email,
		OTPHash: hash,
		ExpiredAt: time.Now().Add(
			time.Minute * time.Duration(u.expiryMinute),
		),
		IPAddress: ip,
		UserAgent: ua,
	}

	return u.otpRepo.Save(ctx, entity)
}

func (u *OTPUsecase) VerifyOTP(
	ctx context.Context,
	email string,
	hash string,
) bool {

	data, err := u.otpRepo.FindActiveByEmail(ctx, email)
	if err != nil {
		return false
	}

	if time.Now().After(data.ExpiredAt) {
		_ = u.otpRepo.DeleteByEmail(ctx, email)
		return false
	}

	if data.OTPHash != hash {
		return false
	}

	_ = u.otpRepo.DeleteByEmail(ctx, email)
	return true
}
