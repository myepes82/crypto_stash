package internal

import (
	"github.com/myepes82/crypto_stash/internal/infrastructure"
	"github.com/myepes82/crypto_stash/internal/services/definitions"
)

type Application struct {
	EncryptionService definitions.ForEncryptingSecretsService
	DecryptionService definitions.ForDecryptingSecrets
}

func NewApplication(
	logger *infrastructure.Logger,
	encryptionService definitions.ForEncryptingSecretsService,
	decryptionService definitions.ForDecryptingSecrets,
) *Application {
	return &Application{
		EncryptionService: encryptionService,
		DecryptionService: decryptionService,
	}
}
