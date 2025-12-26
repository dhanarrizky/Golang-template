package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	ports "github.com/dhanarrizky/Golang-template/internal/ports/others"
)

type hmacTokenVerifier struct {
	secret []byte
}

func NewHMACTokenVerifier(secret string) ports.TokenVerifier {
	return &hmacTokenVerifier{
		secret: []byte(secret),
	}
}

func (h *hmacTokenVerifier) Hash(token string) string {
	mac := hmac.New(sha256.New, h.secret)
	mac.Write([]byte(token))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func (h *hmacTokenVerifier) Compare(hash, token string) bool {
	expected := h.Hash(token)
	return hmac.Equal([]byte(hash), []byte(expected))
}
