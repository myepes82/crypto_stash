package definitions

type ForEncryptingSecretsService interface {
	Encrypt(secret string, key []byte) ([]byte, error)
}
