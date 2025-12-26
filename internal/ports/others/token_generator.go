package others

type TokenGenerator interface {
	Generate() (plain string, hash string, err error)
}
