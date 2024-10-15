package definitions

type ForEncryptingSecretsService interface {
	Execute(secret string, key []byte) ([]byte, error)
}
