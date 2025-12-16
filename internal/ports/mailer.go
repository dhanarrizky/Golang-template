package ports

import "context"

type Mailer interface {
	Send(ctx context.Context, to []string, subject, body string) error
}
