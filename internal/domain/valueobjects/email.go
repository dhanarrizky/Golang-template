// Contoh valueobjects/email.go (opsional, tidak wajib)
package valueobjects

type Email string

func NewEmail(value string) (Email, error) {
	// Simple validation
	if !strings.Contains(value, "@") {
		return "", errors.New("invalid email")
	}
	return Email(value), nil
}