package main

import (
	"github.com/myepes82/crypto_stash/cmd"
	"github.com/myepes82/crypto_stash/internal"
	"github.com/myepes82/crypto_stash/internal/infrastructure"
	"github.com/myepes82/crypto_stash/internal/services"
)

func main() {

	logger := infrastructure.NewLogger()

	//Services
	encryptionService := services.NewEncryptionService(*logger)
	decryptionService := services.NewDecryptionService(*logger)

	//Application
	application := internal.NewApplication(
		logger,
		encryptionService,
		decryptionService,
	)

	console := cmd.NewCmdApplication(logger, application)
	console.Init()
}
