package definitions

type ForDecryptingSecrets interface {
	Execute(secret []byte, key []byte) (string, error)
}
