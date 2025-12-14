package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims used for both access and refresh tokens.
// Keep small and avoid sensitive user data.
type TokenClaims struct {
	UserID string `json:"sub"` // subject is user id
	jwt.RegisteredClaims
}

type Signer struct {
	Secret       []byte
	Issuer       string
	Audience     string
	AccessTTL    time.Duration
	RefreshTTL   time.Duration
}

func NewSigner(secret, issuer, audience string, accessTTL, refreshTTL time.Duration) *Signer {
	return &Signer{
		Secret:     []byte(secret),
		Issuer:     issuer,
		Audience:   audience,
		AccessTTL:  accessTTL,
		RefreshTTL: refreshTTL,
	}
}

func (s *Signer) NewAccessToken(userID string, extraJTI string) (string, *TokenClaims, error) {
	now := time.Now()
	claims := &TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.Issuer,
			Subject:   userID,
			Audience:  jwt.ClaimStrings{s.Audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.AccessTTL)),
			ID:        extraJTI, // optional jti for tracing; can be set blank
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.Secret)
	return signed, claims, err
}

func (s *Signer) NewRefreshToken(userID string, jti string) (string, *TokenClaims, error) {
	now := time.Now()
	claims := &TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.Issuer,
			Subject:   userID,
			Audience:  jwt.ClaimStrings{s.Audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.RefreshTTL)),
			ID:        jti, // jti: important for revocation
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.Secret)
	return signed, claims, err
}

func (s *Signer) ParseToken(tokenStr string) (*TokenClaims, error) {
	claims := &TokenClaims{}
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	_, err := parser.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		return s.Secret, nil
	})

	if err != nil {
		// JWT v5 menggunakan sentinel errors
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token expired")
		}
		return nil, errors.New("token invalid")
	}

	// validate issuer
	if claims.Issuer != s.Issuer {
		return nil, errors.New("invalid issuer")
	}

	// validate audience manually
	validAud := false
	for _, aud := range claims.Audience {
		if aud == s.Audience {
			validAud = true
			break
		}
	}
	if !validAud {
		return nil, errors.New("invalid audience")
	}

	return claims, nil
}
