package valueobjects

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/crypto/argon2"
)

// Password adalah value object yang merepresentasikan password yang sudah di-hash dengan aman.
// Nilai string-nya adalah hashed password (bukan plain text).
type Password string

// PasswordConfig untuk tuning Argon2id (bisa di-inject via config)
type PasswordConfig struct {
	Time    uint32 // iterations
	Memory  uint32 // KiB
	Threads uint8  // parallelism
	KeyLen  uint32 // output length
	SaltLen uint32 // salt length
}

var DefaultPasswordConfig = PasswordConfig{
	Time:    3,         // rekomendasi OWASP 2024-2025
	Memory:  64 * 1024, // 64 MiB
	Threads: 2,
	KeyLen:  32, // 256-bit
	SaltLen: 16, // 128-bit salt
}

// ErrPasswordTooWeak adalah error jika password tidak memenuhi strength policy
var (
	ErrPasswordTooShort  = errors.New("password must be at least 12 characters")
	ErrPasswordTooLong   = errors.New("password cannot exceed 128 characters")
	ErrPasswordNoLower   = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoUpper   = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoDigit   = errors.New("password must contain at least one digit")
	ErrPasswordNoSpecial = errors.New("password must contain at least one special character")
	ErrPasswordCommon    = errors.New("password is too common or easily guessable")
	ErrPasswordPwned     = errors.New("password has been exposed in a data breach") // optional future integration
)

// NewPassword membuat hashed password dari plain text dengan validasi ketat.
// Gunakan ini saat register atau change password.
func NewPassword(plain string, config *PasswordConfig) (Password, error) {
	if config == nil {
		config = &DefaultPasswordConfig
	}

	// 1. Validasi panjang
	if len(plain) < 12 {
		return "", ErrPasswordTooShort
	}
	if len(plain) > 128 {
		return "", ErrPasswordTooLong
	}

	// 2. Validasi komposisi karakter (entropy tinggi)
	if !hasLower(plain) {
		return "", ErrPasswordNoLower
	}
	if !hasUpper(plain) {
		return "", ErrPasswordNoUpper
	}
	if !hasDigit(plain) {
		return "", ErrPasswordNoDigit
	}
	if !hasSpecial(plain) {
		return "", ErrPasswordNoSpecial
	}

	// 3. Optional: cek common password (bisa tambah zxcvbn atau list top 100k)
	// if isCommonPassword(plain) { return "", ErrPasswordCommon }

	// 4. Generate salt acak
	salt := make([]byte, config.SaltLen)
	if _, err := randReader.Read(salt); err != nil {
		return "", err
	}

	// 5. Hash dengan Argon2id (pemenang PHC, rekomendasi OWASP & NIST 2024+)
	hash := argon2.IDKey([]byte(plain), salt, config.Time, config.Memory, config.Threads, config.KeyLen)

	// 6. Format: $argon2id$v=19$m=65536,t=3,p=2$<base64-salt>$<base64-hash>
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	hashed := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		config.Memory, config.Time, config.Threads, encodedSalt, encodedHash)

	return Password(hashed), nil
}

// Verify memeriksa apakah plain password cocok dengan hashed Password
func (p Password) Verify(plain string) (bool, error) {
	// Parse format Argon2
	parts := strings.Split(string(p), "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, errors.New("invalid password format")
	}

	var time, memory uint32
	var threads uint8

	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	expectedHash := argon2.IDKey([]byte(plain), salt, time, memory, threads, uint32(len(hash)))

	return subtle.ConstantTimeCompare(expectedHash, hash) == 1, nil
}

// String mengembalikan hashed value (untuk simpan ke DB)
func (p Password) String() string {
	return string(p)
}

// Helper functions
func hasLower(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

func hasUpper(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func hasDigit(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func hasSpecial(s string) bool {
	for _, r := range s {
		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			return true
		}
	}
	return false
}

// randReader untuk crypto-safe random
var randReader = rand.Reader
