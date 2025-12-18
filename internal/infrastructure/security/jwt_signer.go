package security

import (
	"errors"
	"time"

	"github.com/dhanarrizky/Golang-template/internal/domain/valueobjects"
	ports "github.com/dhanarrizky/Golang-template/internal/ports/auth"
	"github.com/dhanarrizky/Golang-template/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
)

type JWTSigner struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewJWTSigner(
	accessSecret string,
	refreshSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) ports.TokenSigner {
	return &JWTSigner{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (j *JWTSigner) GenerateAccessToken(
	userID string,
	claims map[string]any,
) (string, error) {

	now := time.Now()

	tokenClaims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(j.accessTTL).Unix(),
		"typ": "access",
	}

	for k, v := range claims {
		tokenClaims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	return token.SignedString(j.accessSecret)
}

func (j *JWTSigner) GenerateRefreshToken(userID string) (string, error) {
	now := time.Now()

	tokenID := utils.GenerateID() // wajib unique (jti)

	claims := jwt.MapClaims{
		"sub": userID,
		"jti": tokenID,
		"iat": now.Unix(),
		"exp": now.Add(j.refreshTTL).Unix(),
		"typ": "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(j.refreshSecret)
}

func (j *JWTSigner) VerifyAccessToken(tokenStr string) (*valueobjects.TokenPayload, error) {
	return j.verify(tokenStr, j.accessSecret, "access")
}

func (j *JWTSigner) VerifyRefreshToken(tokenStr string) (*valueobjects.TokenPayload, error) {
	return j.verify(tokenStr, j.refreshSecret, "refresh")
}

func (j *JWTSigner) verify(
	tokenStr string,
	secret []byte,
	expectedType string,
) (*valueobjects.TokenPayload, error) {

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	if claims["typ"] != expectedType {
		return nil, errors.New("invalid token type")
	}

	return &valueobjects.TokenPayload{
		UserID:    claims["sub"].(string),
		TokenID:   claims["jti"].(string),
		ExpiresAt: time.Unix(int64(claims["exp"].(float64)), 0),
	}, nil
}
