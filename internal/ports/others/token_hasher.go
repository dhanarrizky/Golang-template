package others

type TokenHasher interface {
	Hash(token string) string
}
