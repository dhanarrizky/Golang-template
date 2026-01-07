// package email

// import (
// 	"net/smtp"

// 	ports "github.com/dhanarrizky/Golang-template/internal/ports/email"
// )

// type SMTPEmailService struct {
// 	host string
// 	port string
// 	from string
// 	auth smtp.Auth
// }

// func NewSMTPEmailService(host, port, from, pw string) ports.EmailSender {
// 	return &SMTPEmailService{
// 		host: host,
// 		port: port,
// 		from: from,
// 		auth: smtp.PlainAuth("", from, pw, host),
// 	}
// }

// func (s *SMTPEmailService) SendOTP(toEmail, otp string) error {
// 	// render template di sini
// 	return smtp.SendMail(
// 		s.host+":"+s.port,
// 		s.auth,
// 		s.from,
// 		[]string{toEmail},
// 		[]byte("Subject: OTP\r\n\r\nYour OTP: "+otp),
// 	)
// }

// internal/adapters/email/smtp_sender.go
package email

import (
	"fmt"
	"net/smtp"
)

type SMTPSender struct {
	host string
	port string
	user string
	pass string
	from string
}

func NewSMTPSender(host, port, user, pass, from string) *SMTPSender {
	return &SMTPSender{host, port, user, pass, from}
}

func (s *SMTPSender) SendOTP(to, otp string) error {
	body := fmt.Sprintf(
		"Subject: Your OTP Code\n\nYour OTP is: %s",
		otp,
	)

	addr := s.host + ":" + s.port
	auth := smtp.PlainAuth("", s.user, s.pass, s.host)

	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(body))
}

func (s *SMTPSender) SendResetPassword(to, otp string) error {
	return s.SendOTP(to, otp)
}
