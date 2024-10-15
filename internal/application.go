package internal

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/myepes82/crypto_stash/internal/infrastructure"
	"github.com/myepes82/crypto_stash/internal/models"
	"github.com/myepes82/crypto_stash/internal/services/definitions"
	"os"
)

type Application struct {
	//Secret key
	secrets *models.Secrets

	//Logger
	Logger *infrastructure.Logger

	//Services
	encryptionService        definitions.ForEncryptingSecretsService
	decryptionService        definitions.ForDecryptingSecrets
	creatingSecretKeyService definitions.ForCreatingSecretKey
}

func NewApplication(
	logger *infrastructure.Logger,
	encryptionService definitions.ForEncryptingSecretsService,
	decryptionService definitions.ForDecryptingSecrets,
	creatingSecretKeyService definitions.ForCreatingSecretKey,
) *Application {
	return &Application{
		Logger:                   logger,
		encryptionService:        encryptionService,
		decryptionService:        decryptionService,
		creatingSecretKeyService: creatingSecretKeyService,
	}
}

func (app *Application) DecryptSecret(secret []byte, key []byte) (string, error) {
	return app.decryptionService.Execute(secret, key)
}

func (app *Application) CreateSecretKey() ([]byte, error) {
	return app.creatingSecretKeyService.Execute()
}

func (app *Application) LoadSecrets(secrets *models.Secrets) {
	if secrets != nil {
		app.secrets = secrets
	} else {
		app.secrets = models.NewSecrets()
	}
}

func (app *Application) AddSecret(name string, secret string, secretKey []byte) error {
	key, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	encryptedSecret, err := app.encryptionService.Execute(secret, secretKey)
	if err != nil {
		return err
	}

	app.secrets.AddSecret(name, key.String())

	go func() {
		if err := app.generateSecretFile(key.String(), encryptedSecret); err != nil {
			fmt.Printf("Error generating secret file: %v\n", err)
		}
	}()

	return nil
}

func (app *Application) GetSecretsRaw() models.Secrets {
	if app.secrets != nil {
		return *app.secrets
	}
	return models.Secrets{}
}

func (app *Application) GetSecret(key string, secretKey []byte) (string, error) {
	secretFound, ok := app.secrets.GetSecret(key)
	if !ok {
		return "", errors.New("secret not found")
	}

	secretFilePath := fmt.Sprintf("./secrets/%s.txt", secretFound)

	_, err := os.Stat(secretFilePath)
	if err != nil && os.IsNotExist(err) {
		return "", errors.New("secret not found")
	}

	data, err := os.ReadFile(secretFilePath)

	if err != nil {
		return "", err
	}

	return app.decryptionService.Execute(data, secretKey)
}

func (app *Application) generateSecretFile(name string, secret []byte) error {
	fileName := name + ".txt"
	filePath := fmt.Sprintf("./secrets/%s", fileName)

	_, err := os.Stat("secrets")
	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir("secrets", 0755)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(filePath, secret, 0644)

	fmt.Printf("generated secret file: %s\n", filePath)
	return nil
}
