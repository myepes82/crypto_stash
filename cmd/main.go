package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
)

const (
	Reset           = "\033[0m"
	Red             = "\033[31m"
	Green           = "\033[32m"
	Yellow          = "\033[33m"
	SecretsFileName = "secrets.json"
)

type SecretsStruct struct {
	Changelog time.Time         `json:"changelog"`
	Secrets   map[string]string `json:"secrets"`
}

func main() {
	var secretsData SecretsStruct
	var secretKey string
	welcome_logo := `    
                      | |           | |          | |    
  ___ _ __ _   _ _ __ | |_ ___   ___| |_ __ _ ___| |__  
 / __| '__| | | | '_ \| __/ _ \ / __| __/ _ \ /__| '_  \ 
| (__| |  | |_| | |_) |  | (_) |\__ \  | (_| \__ \ | | |
 \___|_|   \__, | .__/ \__\___/ |___/\__\__,_|___/_| |_|
            __/ | |                               
           |___/|_|                            
	`

	fmt.Println(Red + welcome_logo + Reset)
	fmt.Println(
		Yellow +
			"Welcome to crypto_stash an open tool for keep secrets secured" +
			Reset)

	keyPathFlag := flag.String("key", "", "Key for signing secrets")
	flag.Parse()

	loadSecretKeyFile(&secretKey, *keyPathFlag)

	fmt.Println(secretKey)
	err := checkExistenceSecrets()

	if err != nil {
		fmt.Println(Red + "Error in file checking process." + Reset)
	}

	err = loadSecretsInfo(&secretsData)

	if err != nil {
		fmt.Println(Red + "Error while loading secrets info process." + Reset)
	}

	cmd_options := []string{
		"1) List secrets",
		"2) Add new secret",
		"3) Get one secret",
	}

	prompt := promptui.Select{
		Label: "Choose what you want to do",
		Items: cmd_options,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf(Red+"An unexpected error has ocurred: %v"+Reset, err)
		os.Exit(0)
	}

	chosenOption, err := strconv.Atoi(result[:1])
	switch chosenOption {
	case 1:
		listSecrets(secretsData)
	default:
		fmt.Println("ToDo")
	}
	fmt.Println(chosenOption)
}

func loadSecretKeyFile(secretKey *string, secretKeyFilePath string) error {
	_, err := os.Stat(secretKeyFilePath)
	if !os.IsNotExist(err) {
		data, err := os.ReadFile(secretKeyFilePath)
		if err != nil {
			println(Red+"Error reading secret key file: %v"+Reset, err)
		}
		*secretKey = string(data)
		return nil
	}
	fmt.Println("No secret key file were found")
	cmd_options := []string{
		"Yes",
		"No",
	}

	prompt := promptui.Select{
		Label: "Want to generate a secret key?",
		Items: cmd_options,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf(Red+"An unexpected error has ocurred: %v"+Reset, err)
		return err
	}

	if result == "No" {
		os.Exit(0)
	}
	_, err = os.Create("key.txt")

	if err != nil {
		fmt.Printf(Red+"An unexpected error has ocurred: %v"+Reset, err)
		os.Exit(0)
	}
	return nil
}

func listSecrets(secretsData SecretsStruct) {
	fmt.Println(Green + "============ Listing secrets ============" + Reset)
	for key, value := range secretsData.Secrets {
		fmt.Printf(Yellow+"Name: "+Reset+"%s "+Yellow+"Secret: "+Reset+"%s", key, value)
	}
}

func loadSecretsInfo(secretsData *SecretsStruct) error {
	data, err := os.ReadFile(SecretsFileName)

	if err != nil {
		fmt.Printf(Red+"An error has ocurred while reading secrets file: %v"+Reset, err)
		return err
	}

	err = json.Unmarshal(data, &secretsData)

	if err != nil {
		fmt.Printf(Red+"An error has ocurred while unmarshal secrets file info: %v"+Reset, err)
		return err
	}
	return nil
}
func createInitialSecretsFile() error {
	secretsTemplate := SecretsStruct{
		Changelog: time.Now(),
		Secrets:   make(map[string]string),
	}

	jsonData, err := json.MarshalIndent(secretsTemplate, "", "    ")

	if err != nil {
		fmt.Printf(Red+"Error marshalling JSON: %v"+Reset, err)
		return err
	}

	err = os.WriteFile(SecretsFileName, jsonData, 0644)

	if err != nil {
		fmt.Printf(Red+"Error creating secrets.json file: %v"+Reset, err)
		return err
	}
	fmt.Println("The file was created")
	return nil
}

func checkExistenceSecrets() error {
	_, err := os.Stat(SecretsFileName)
	if !os.IsNotExist(err) {
		return nil
	}
	fmt.Println(Green + "Secrets file not found, creating..." + Reset)

	return createInitialSecretsFile()
}
