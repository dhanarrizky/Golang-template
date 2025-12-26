package users

// PasswordHasher mendefinisikan kontrak untuk hashing dan verifikasi password
type PasswordHasher interface {
	// HashPassword menghasilkan hash password dengan pepper versi saat ini
	// HashPassword(password []byte) (hashed string, pepperVersion int, err error)
	HashPassword(password []byte) (string, error)

	// VerifyPassword memverifikasi password terhadap hash yang sudah ada
	// Mengembalikan (match bool, shouldRehash bool, err error)
	VerifyPassword(password []byte, hashed string) (bool, bool, error)
}
