package auth

import "github.com/dhanarrizky/Golang-template/internal/ports"

type SendOTPUsecase struct {
	email ports.EmailSender
}

func NewSendOTPUsecase(email ports.EmailSender) *SendOTPUsecase {
	return &SendOTPUsecase{email: email}
}

func (uc *SendOTPUsecase) Execute(email string, otp string) error {
	return uc.email.SendOTP(email, otp)
}
