package others

type TokenVerifier interface {
	Hash(token string) string
	Compare(hash, token string) bool
}
