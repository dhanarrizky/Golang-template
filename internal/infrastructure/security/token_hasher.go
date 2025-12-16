package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"github.com/dhanarrizky/Golang-template/internal/ports"
)

type HMACTokenHasher struct {
	secret []byte
}

func NewHMACTokenHasher(secret string) ports.TokenHasher {
	return &HMACTokenHasher{
		secret: []byte(secret),
	}
}

func (h *HMACTokenHasher) Hash(token string) string {
	mac := hmac.New(sha256.New, h.secret)
	mac.Write([]byte(token))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
