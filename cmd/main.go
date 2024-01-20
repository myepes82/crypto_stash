package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
)

type SecretsStruct struct {
	Changelog time.Time         `json:"changelog"`
	Secrets   map[string]string `json:"secrets"`
}

func (s *SecretsStruct) UpdateChangelog() {
	s.Changelog = time.Now()
}
func (s *SecretsStruct) AddSecret(key, value string) {
	if s.Secrets == nil {
		s.Secrets = make(map[string]string)
	}
	s.Secrets[key] = value
}

func (s *SecretsStruct) GetSecret(key string) (string, bool) {
	value, ok := s.Secrets[key]
	return value, ok
}

const (
	Reset           = "\033[0m"
	Red             = "\033[31m"
	Green           = "\033[32m"
	Yellow          = "\033[33m"
	SecretsFileName = "secrets.json"
)

func main() {
	var secretsData SecretsStruct
	var secretKey []byte

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

	keyPathFlag := flag.String("key", "secret_key.txt", "Key for signing secrets")
	flag.Parse()

	err := loadSecretKeyFile(&secretKey, *keyPathFlag)

	if err != nil {
		handleUnknownError(err)
	}

	fmt.Println(Green + "Secret key loaded successful" + Reset)

	err = checkExistenceSecrets()

	if err != nil {
		fmt.Println(Red + "Error in secrets file checking process." + Reset)
		os.Exit(0)
	}

	err = loadSecretsInfo(&secretsData)

	if err != nil {
		fmt.Println(Red + "Error while loading secrets info process." + Reset)
		os.Exit(0)
	}

	cmdOptions := []string{
		"1) List secrets",
		"2) Add new secret",
		"3) Get one secret",
	}

	prompt := promptui.Select{
		Label: "Choose what you want to do",
		Items: cmdOptions,
	}

	for {
		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf(Red+"An unexpected error has ocurred: %v"+Reset, err)
			break
		}
		chosenOption, err := strconv.Atoi(result[:1])
		processCryptoAction(chosenOption, &secretsData, secretKey)
	}

}
func processCryptoAction(chosenOption int, secretsData *SecretsStruct, secretKey []byte) {
	switch chosenOption {
	case 1:
		listSecrets(*secretsData)
	case 2:
		addNewSecret(secretsData, secretKey)
	default:
		fmt.Println("ToDo")
	}
}
func handleUnknownError(err error) {
	fmt.Printf(Red+"An unexpected error has occurred: %v"+Reset, err)
	os.Exit(1)
}
func loadSecretKeyFile(secretKey *[]byte, secretKeyFilePath string) error {
	_, err := os.Stat(secretKeyFilePath)
	if !os.IsNotExist(err) {
		data, err := os.ReadFile(secretKeyFilePath)
		if err != nil {
			handleUnknownError(err)
		}
		*secretKey = data
		return nil
	}
	fmt.Println("No secret key file were found")
	cmdOptions := []string{
		"Yes",
		"No",
	}

	prompt := promptui.Select{
		Label: "Want to generate a secret key?",
		Items: cmdOptions,
	}

	_, result, err := prompt.Run()

	if err != nil {
		handleUnknownError(err)
	}

	if result == "No" {
		os.Exit(0)
	}
	key, err := generateSecretKey()

	if err != nil {
		handleUnknownError(err)
	}

	generateSecretKeyFile(key)

	*secretKey = key
	return nil
}
func generateSecretKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return []byte{}, err
	}
	return key, nil
}
func generateSecretKeyFile(key []byte) {
	err := os.WriteFile("secret_key.txt", key, 0644)
	if err != nil {
		handleUnknownError(err)
	}
	println(Green + "secret_key.txt file was created" + Reset)
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
func updateSecretsFile(secrets *SecretsStruct) error {
	data, err := json.MarshalIndent(secrets, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile(SecretsFileName, data, 0644)
	if err != nil {
		return err
	}
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
func listSecrets(secretsData SecretsStruct) {
	fmt.Println(Green + "================= Listing secrets ================" + Reset)
	for key, value := range secretsData.Secrets {
		fmt.Printf(Yellow+"Name: "+Reset+"%s "+Yellow+"Secret: "+Reset+"%s \n", key, value)
	}
	fmt.Println(Green + "========================||========================" + Reset)
}
func encryptSecret(secret string, secretKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ciphertext := aesGCM.Seal(nil, nonce, []byte(secret), nil)

	ciphertext = append(nonce, ciphertext...)

	return ciphertext, nil
}
func generateSecretTextFile(secret string, secretKey []byte, secretUUid string) error {

	encryptedSecret, err := encryptSecret(secret, secretKey)
	if err != nil {
		handleUnknownError(err)
	}
	fileName := fmt.Sprintf("./secrets/%s.txt", secretUUid)

	err = os.WriteFile(fileName, encryptedSecret, 0644)
	if err != nil {
		handleUnknownError(err)
	}
	println(Green + fileName + "file was created" + Reset)
	return nil
}
func decryptSecret(secret []byte, secretKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	nonce := secret[:aes.BlockSize]
	secret = secret[aes.BlockSize:]

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesGCM.Open(nil, nonce, secret, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
func addNewSecret(secrets *SecretsStruct, secretKey []byte) {
	fmt.Println(Red + "Please fill the following fields" + Reset)
	namePrompt := promptui.Prompt{
		Label: "Name: ",
	}
	secretPrompt := promptui.Prompt{
		Label: "Secret: ",
	}

	name, err := namePrompt.Run()

	if err != nil {
		handleUnknownError(err)
	}
	secret, err := secretPrompt.Run()

	if err != nil {
		handleUnknownError(err)
	}

	secretUUid, err := uuid.FromBytes([]byte(secret))
	if err != nil {
		handleUnknownError(err)
	}
	secretUUidString := secretUUid.String()

	secrets.AddSecret(name, secretUUidString)
	secrets.UpdateChangelog()

	err = updateSecretsFile(secrets)
	if err != nil {
		handleUnknownError(err)
	}

	go generateSecretTextFile(secret, secretKey, secretUUidString)
	fmt.Println(Green + "New secret text file generated" + Reset)
}
