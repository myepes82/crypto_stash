package services

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/myepes82/crypto_stash/internal/infrastructure"
)

type DecryptionService struct {
	logger infrastructure.Logger
}

func NewDecryptionService(logger infrastructure.Logger) *DecryptionService {
	return &DecryptionService{
		logger: logger,
	}
}

func (service *DecryptionService) Execute(secret []byte, key []byte) (string, error) {
	service.logger.LogDebug("")
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := secret[:aes.BlockSize]
	secret = secret[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(secret, secret)

	plaintext := pkcs7Unpad(secret)

	return string(plaintext), nil
}

func pkcs7Unpad(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}
