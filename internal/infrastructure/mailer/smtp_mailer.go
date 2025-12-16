package mailer

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strconv"
	"strings"
)

type SMTPMailer struct {
	host        string
	port        string
	username    string
	password    string
	from        string
	fromAddress string
	auth        smtp.Auth
}

func NewSMTPMailer(
	host, port, username, password, fromName, fromAddress string,
) (*SMTPMailer, error) {

	if host == "" || port == "" || username == "" || password == "" {
		return nil, fmt.Errorf("smtp configuration incomplete")
	}

	if fromName == "" {
		fromName = "App"
	}

	auth := smtp.PlainAuth("", username, password, host)

	return &SMTPMailer{
		host:        host,
		port:        port,
		username:    username,
		password:    password,
		from:        fmt.Sprintf("%s <%s>", fromName, fromAddress),
		fromAddress: fromAddress,
		auth:        auth,
	}, nil
}

func (m *SMTPMailer) Send(
	ctx context.Context,
	to []string,
	subject, body string,
) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	msg := []byte(fmt.Sprintf(
		"To: %s\r\nFrom: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		strings.Join(to, ","),
		m.from,
		subject,
		body,
	))

	addr := m.host + ":" + m.port
	portInt, _ := strconv.Atoi(m.port)

	// SSL 465
	if portInt == 465 {
		conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: m.host})
		if err != nil {
			return err
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, m.host)
		if err != nil {
			return err
		}
		defer client.Close()

		if err = client.Auth(m.auth); err != nil {
			return err
		}

		if err = client.Mail(m.fromAddress); err != nil {
			return err
		}

		for _, r := range to {
			if err = client.Rcpt(r); err != nil {
				return err
			}
		}

		w, err := client.Data()
		if err != nil {
			return err
		}

		if _, err = w.Write(msg); err != nil {
			return err
		}

		return w.Close()
	}

	// STARTTLS (587)
	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); ok {
		if err = client.StartTLS(&tls.Config{ServerName: m.host}); err != nil {
			return err
		}
	}

	if err = client.Auth(m.auth); err != nil {
		return err
	}

	if err = client.Mail(m.fromAddress); err != nil {
		return err
	}

	for _, r := range to {
		if err = client.Rcpt(r); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	return w.Close()
}
