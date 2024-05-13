package definitions

type ForDecryptingSecrets interface {
	Decrypt(secret []byte, key []byte) (string, error)
}
