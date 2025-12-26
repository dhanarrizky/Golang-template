package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"

	ports "github.com/dhanarrizky/Golang-template/internal/ports/others"
)

type aesCodec struct {
	key []byte
}

func NewPublicIDCodecFromBase64(keyBase64 string) (ports.PublicIDCodec, error) {
	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return nil, errors.New("invalid base64 AES key")
	}

	switch len(key) {
	case 16, 24, 32:
		// valid
	default:
		return nil, errors.New("AES key must be 16, 24, or 32 bytes")
	}

	return &aesCodec{key: key}, nil
}

func (a *aesCodec) Encode(id uint64) (string, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, id)

	encrypted := gcm.Seal(nonce, nonce, buf, nil)
	return base64.RawURLEncoding.EncodeToString(encrypted), nil
}

func (a *aesCodec) Decode(publicID string) (uint64, error) {
	data, err := base64.RawURLEncoding.DecodeString(publicID)
	if err != nil {
		return 0, err
	}

	block, err := aes.NewCipher(a.key)
	if err != nil {
		return 0, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return 0, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return 0, errors.New("invalid data")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(plain), nil
}
