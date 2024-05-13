package services

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/myepes82/crypto_stash/internal/infrastructure"
)

type EncryptionService struct {
	logger infrastructure.Logger
}

func NewEncryptionService(logger infrastructure.Logger) *EncryptionService {
	return &EncryptionService{
		logger: logger,
	}
}

func (service *EncryptionService) Encrypt(secret string, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	paddedSecret := pkcs7Pad([]byte(secret), aes.BlockSize)

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)

	ciphertext := make([]byte, len(paddedSecret))
	mode.CryptBlocks(ciphertext, paddedSecret)

	result := append(iv, ciphertext...)

	return result, nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}
