package others

// type TokenHasher interface {
// 	Hash(token string) string
// }

type PublicIDCodec interface {
	Encode(id uint64) (string, error)
	Decode(publicID string) (uint64, error)
}
