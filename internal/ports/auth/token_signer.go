package auth

import "github.com/dhanarrizky/Golang-template/internal/domain/valueobjects"

// token_signer.go
type TokenSigner interface {
	GenerateAccessToken(userID string, claims map[string]any) (string, error)
	GenerateRefreshToken(userID string) (string, error)

	VerifyAccessToken(token string) (*valueobjects.TokenPayload, error)
	VerifyRefreshToken(token string) (*valueobjects.TokenPayload, error)
}
