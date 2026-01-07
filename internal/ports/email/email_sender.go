package email

// type EmailSender interface {
// 	SendOTP(to string, otp string) error
// }

// internal/ports/email/sender.go

type EmailSender interface {
	SendOTP(to string, otp string) error
	SendResetPassword(to string, otp string) error
}
