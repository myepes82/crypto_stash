package cmd

import (
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/myepes82/crypto_stash/internal"
	"github.com/myepes82/crypto_stash/internal/infrastructure"
)

const (

	//Labels
	wantToGenerateSecretLabel = "do u want to generate a secret key?"
	cmdMenuOptionsLabel       = "choose what you want to do:"

	//Message
	welcomeMessage                         = "welcome to crypto_stash an open tool for keep secrets secured."
	closingAppMessage                      = "closing app."
	exitingApplicationMessage              = "exiting application"
	noActionSelectedMessage                = "no action selected."
	noSecretsKeyFileWillBeGeneratedMessage = "no secrets key file will be generated."
	noSecretsKeyFileFoundMessage           = "no secret key file were found."
)

var (
	yerOrNotCmdOptions = []string{
		"Yes",
		"No",
	}
	cmdMenuOptions = []string{
		"1) List secrets",
		"2) Add new secret",
		"3) Get one secret",
		"0) Exit",
	}
)

type Cmd struct {
	logger    *infrastructure.Logger
	app       *internal.Application
	secretKey []byte
}

func NewCmdApplication(
	logger *infrastructure.Logger,
	app *internal.Application) *Cmd {
	return &Cmd{
		logger: logger,
		app:    app,
	}
}

func (cmd *Cmd) showInitialMessage() {
	welcomeLogoMessage :=
		`    
						| |           | |          | |    
	___ _ __ _   _ _ __ | |_ ___   ___| |_ __ _ ___| |__  
   / __| '__| | | | '_ \| __/ _ \ / __| __/ _ \ /__| '_  \ 
  | (__| |  | |_| | |_) |  | (_) |\__ \  | (_| \__ \ | | |
   \___|_|   \__, | .__/ \__\___/ |___/\__\__,_|___/_| |_|
              __/ | |                               
			 |___/|_|                            
	`
	cmd.logger.LogDebug(welcomeLogoMessage)
	cmd.logger.LogWarm(welcomeMessage)
}

func (cmd *Cmd) createSecretsKeyFile() error {
	if err := os.WriteFile("secret_key.txt", cmd.secretKey, 0644); err != nil {
		return err
	}
	cmd.logger.LogDebug("secret_key.txt file was created")
	return nil
}

func (cmd *Cmd) createSecretsKey() error {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return err
	}
	cmd.secretKey = key
	return nil
}

func (cmd *Cmd) loadSecretsKeyFile() error {
	keyPathFlag := *flag.String("key", "secret_key.txt", "Key for signing secrets")
	flag.Parse()

	_, err := os.Stat(keyPathFlag)
	if !os.IsNotExist(err) {
		data, err := os.ReadFile(keyPathFlag)
		if err != nil {
			cmd.handleUnknownError(err)
		}
		*&cmd.secretKey = data
	}

	cmd.logger.LogWarm(noSecretsKeyFileFoundMessage)
	return errors.New(noSecretsKeyFileFoundMessage)
}

func (cmd *Cmd) processNoExitingSecretsFile() {
	prompt := promptui.Select{
		Label: wantToGenerateSecretLabel,
		Items: yerOrNotCmdOptions,
	}
	_, result, err := prompt.Run()
	if err != nil {
		cmd.handleUnknownError(err)
	}
	if result == "No" {
		cmd.logger.LogDebug(noSecretsKeyFileWillBeGeneratedMessage)
		cmd.logger.LogDebug(closingAppMessage)
		os.Exit(0)
	}
	if err := cmd.createSecretsKey(); err != nil {
		cmd.handleUnknownError(err)
	}
	if err := cmd.createSecretsKeyFile(); err != nil {
		cmd.handleUnknownError(err)
	}
}

func (cmd *Cmd) proccessAction(option int) bool {
	switch option {
	case 1:
		break
	case 2:
		break
	case 3:
		break
	case 0:
		cmd.logger.LogDebug(exitingApplicationMessage)
		return false
	default:
		cmd.logger.LogWarm(noActionSelectedMessage)
		return true
	}
	return true
}

func (cmd *Cmd) Init() {
	cmd.showInitialMessage()
	if err := cmd.loadSecretsKeyFile(); err != nil {
		cmd.processNoExitingSecretsFile()
	}
	prompt := promptui.Select{
		Label: cmdMenuOptionsLabel,
		Items: cmdMenuOptions,
	}

	for {
		_, result, err := prompt.Run()
		if err != nil {
			cmd.handleUnknownError(err)
			break
		}
		chosenOption, err := strconv.Atoi(result[:1])
		if !processAction(chosenOption, &secretsData, secretKey) {
			break
		}
	}
}

func (cmd *Cmd) handleUnknownError(err error) {
	cmd.app.Logger.LogError(fmt.Sprintf("unexpected error has occurred: %v", err))
	os.Exit(1)
}
