package definitions

type ForCreatingSecretKey interface {
	Execute() ([]byte, error)
}
