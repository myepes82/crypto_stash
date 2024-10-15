package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/myepes82/crypto_stash/cmd/errors"
	"github.com/myepes82/crypto_stash/internal"
	"github.com/myepes82/crypto_stash/internal/models"
	"os"
	"strconv"
)

const (

	//Labels
	wantToGenerateSecretLabel = "do u want to generate a secret key?"
	cmdMenuOptionsLabel       = "choose what you want to do:"
	secretsFileName           = "secrets.json"

	//Message
	welcomeMessage                         = "welcome to crypto_stash an open tool for keep secrets secured."
	closingAppMessage                      = "closing app."
	exitingApplicationMessage              = "exiting application"
	noActionSelectedMessage                = "no action selected."
	noSecretsKeyFileWillBeGeneratedMessage = "no secrets key file will be generated."

	//Colors
	Reset  = "\033[0m"
	Yellow = "\033[33m"
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
	app       *internal.Application
	secretKey []byte
}

func NewCmdApplication(
	app *internal.Application) *Cmd {
	return &Cmd{
		app: app,
	}
}

func (cmd *Cmd) showInitialMessage() {
	welcomeLogoMessage := `    
                      | |           | |          | |    
  ___ _ __ _   _ _ __ | |_ ___   ___| |_ __ _ ___| |__  
 / __| '__| | | | '_ \| __/ _ \ / __| __/ _ \ /__| '_  \ 
| (__| |  | |_| | |_) |  | (_) |\__ \  | (_| \__ \ | | |
 \___|_|   \__, | .__/ \__\___/ |___/\__\__,_|___/_| |_|
            __/ | |                               
           |___/|_|                            
	`
	cmd.app.Logger.LogDebug(welcomeLogoMessage)
	cmd.app.Logger.LogWarm(welcomeMessage)
}

func (cmd *Cmd) createSecretsKeyFile() error {
	if err := os.WriteFile("secret_key.txt", cmd.secretKey, 0644); err != nil {
		return err
	}
	cmd.app.Logger.LogDebug("secret_key.txt file was created")
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
		return nil
	}

	cmd.app.Logger.LogError(errors.NoSecretKeyFileFoundError)
	return errors.NoSecretKeyFileFoundError
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
		cmd.app.Logger.LogDebug(noSecretsKeyFileWillBeGeneratedMessage)
		cmd.app.Logger.LogDebug(closingAppMessage)
		os.Exit(0)
	}

	secretKey, err := cmd.app.CreateSecretKey()
	if err != nil {
		cmd.handleUnknownError(err)
	}
	cmd.secretKey = secretKey

	if err := cmd.createSecretsKeyFile(); err != nil {
		cmd.handleUnknownError(err)
	}
}

func (cmd *Cmd) processAction(option int) bool {
	switch option {
	case 1:
		cmd.listSecrets()
		return true
	case 2:
		cmd.addSecret()
		return true
	case 3:
		cmd.getOneSecret()
		return true
	case 0:
		cmd.app.Logger.LogDebug(exitingApplicationMessage)
		return false
	default:
		cmd.app.Logger.LogWarm(noActionSelectedMessage)
		return false
	}
}

func (cmd *Cmd) handleUnknownError(err error) {
	cmd.app.Logger.LogError(errors.NewUnknownError(err.Error()))
	os.Exit(1)
}

func (cmd *Cmd) checkExistenceSecrets() error {
	_, err := os.Stat(secretsFileName)
	if !os.IsNotExist(err) {
		return nil
	}
	cmd.app.Logger.LogDebug("Secrets file not found, creating...")

	return cmd.createInitialSecretsFile()
}

func (cmd *Cmd) createInitialSecretsFile() error {
	secretsTemplate := cmd.app.GetSecretsRaw()

	jsonData, err := json.MarshalIndent(secretsTemplate, "", "    ")

	if err != nil {
		parsedError := errors.NewUnknownError(err.Error())
		cmd.app.Logger.LogError(parsedError)
		return parsedError
	}

	err = os.WriteFile(secretsFileName, jsonData, 0644)

	if err != nil {
		parsedError := errors.NewUnknownError(err.Error())
		cmd.app.Logger.LogError(parsedError)
		return parsedError
	}
	cmd.app.Logger.LogDebug("secrets file generated")
	return nil
}

func (cmd *Cmd) loadSecretsInfo() (*models.Secrets, error) {
	data, err := os.ReadFile(secretsFileName)
	sec := &models.Secrets{}

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &sec)

	if err != nil {
		cmd.app.Logger.LogError(fmt.Errorf("error parsing secrets file: %s", err))
		return nil, err
	}
	return sec, nil
}

func (cmd *Cmd) updateSecretsFile() error {
	secrets := cmd.app.GetSecretsRaw()
	data, err := json.MarshalIndent(secrets, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile(secretsFileName, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Specific function

func (cmd *Cmd) addSecret() {
	cmd.app.Logger.LogDebug("Please fill the following fields")
	namePrompt := promptui.Prompt{
		Label: "Name: ",
	}
	secretPrompt := promptui.Prompt{
		Label: "Secret: ",
	}

	name, err := namePrompt.Run()

	if err != nil {
		cmd.handleUnknownError(err)
	}
	secret, err := secretPrompt.Run()

	if err != nil {
		cmd.handleUnknownError(err)
	}

	if err := cmd.app.AddSecret(name, secret, cmd.secretKey); err != nil {
		cmd.handleUnknownError(err)
	}

	go func() {
		if err := cmd.updateSecretsFile(); err != nil {
			cmd.app.Logger.LogError(err)
			return
		}
	}()
}

func (cmd *Cmd) listSecrets() {
	cmd.app.Logger.LogDebug("================= Listing secrets ================")
	for key, value := range cmd.app.GetSecretsRaw().Content {
		fmt.Printf(Yellow+"Name: "+Reset+"%s "+Yellow+"Secret: "+Reset+"%s \n", key, value)
	}
	cmd.app.Logger.LogDebug("========================||========================")
}

func (cmd *Cmd) getOneSecret() {
	var secretList []string
	for key := range cmd.app.GetSecretsRaw().Content {
		secretList = append(secretList, key)
	}
	promptList := promptui.Select{
		Label: "Secrets",
		Items: secretList,
	}

	_, chosenSecret, err := promptList.Run()

	if err != nil {
		cmd.handleUnknownError(err)
	}

	decryptedData, err := cmd.app.GetSecret(chosenSecret, cmd.secretKey)
	if err != nil {
		cmd.handleUnknownError(err)
	}

	fmt.Println(Yellow + "Secret🔐: " + Reset + decryptedData)
}

func (cmd *Cmd) Init() {
	cmd.showInitialMessage()
	err := cmd.checkExistenceSecrets()
	if err != nil {
		cmd.handleUnknownError(err)
	}

	sec, err := cmd.loadSecretsInfo()

	if err != nil {
		cmd.handleUnknownError(err)
	}

	cmd.app.LoadSecrets(sec)

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
		if !cmd.processAction(chosenOption) {
			break
		}
	}
}
