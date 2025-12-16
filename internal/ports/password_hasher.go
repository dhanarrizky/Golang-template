package ports

type PasswordHasher interface {
	Hash(plain string) (hashed string, err error)
	Compare(plain, hashed string) (match bool, err error)
}
