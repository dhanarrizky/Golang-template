// package email

// import "github.com/dhanarrizky/Golang-template/internal/domain/valueobjects"

// type OTPRepository interface {
// 	Save(email string, otp valueobjects.OTP) error
// 	Find(email string) (valueobjects.OTP, bool)
// 	Delete(email string) error
// }

package email

import (
	"context"

	domain "github.com/dhanarrizky/Golang-template/internal/domain/entities/auth"
)

type EmailOTPRepository interface {
	Save(ctx context.Context, otp *domain.EmailOTP) error
	FindActiveByEmail(
		ctx context.Context,
		email string,
	) (*domain.EmailOTP, error)
	DeleteByEmail(ctx context.Context, email string) error
}
