package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/myepes82/crypto_stash/internal/infrastructure"
)

type CreatingSecretKeyService struct {
	logger infrastructure.Logger
}

func NewCreatingSecretKeyService(logger infrastructure.Logger) *CreatingSecretKeyService {
	return &CreatingSecretKeyService{logger: logger}
}

func (service *CreatingSecretKeyService) Execute() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		service.logger.LogError(errors.New(fmt.Sprintf("error generating secret key: %s", err.Error())))
		return nil, err
	}
	return key, nil
}
