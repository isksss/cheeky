package commands

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/google/subcommands"
	"github.com/isksss/cheeky/config"
	"github.com/mikesmitty/edkey"
	"golang.org/x/crypto/ssh"
)

type NewCmd struct{}

func (*NewCmd) Name() string {
	return "new"
}

func (*NewCmd) Synopsis() string {
	return "Create a new key pair"
}

func (*NewCmd) Usage() string {
	return "new <name> <email>"
}

func (*NewCmd) SetFlags(f *flag.FlagSet) {}

func (*NewCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() < 2 {
		log.Printf("You must provide a name for the new directory and an email.")
		return subcommands.ExitFailure
	}
	dirname := f.Arg(0)
	email := f.Arg(1)
	newDir := filepath.Join(config.GetKeysDir(), dirname)

	// Check and create directory
	if err := checkAndCreateDir(newDir); err != nil {
		log.Printf("%v", err)
		return subcommands.ExitFailure
	}

	// Create config file
	if err := createConfigAndKeyFiles(newDir, email); err != nil {
		log.Printf("%v", err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func checkAndCreateDir(newDir string) error {
	// Check if directory already exists
	if directoryExists(newDir) {
		return fmt.Errorf("Directory '%s' already exists", newDir)
	}
	// Create directory
	config.CreateDirIfNotExist(newDir)
	return nil
}

func createConfigAndKeyFiles(newDir string, email string) error {
	// Create config file
	if err := CreateConfigFile(newDir, email); err != nil {
		return err
	}

	pubKeyPath := filepath.Join(newDir, "id_ed25519.pub")
	privKeyPath := filepath.Join(newDir, "id_ed25519")

	pubKey, privKey, err := GenerateKeyPair()
	if err != nil {
		return err
	}

	if err := WriteKeysToFile(pubKeyPath, privKeyPath, pubKey, privKey); err != nil {
		return err
	}

	fmt.Printf("Public key saved to: %s\n", pubKeyPath)
	fmt.Printf("Private key saved to: %s\n", privKeyPath)
	return nil
}

func CreateConfigFile(dirPath string, email string) error {
	configFilePath := filepath.Join(dirPath, "config.json")

	// Check if config file already exists
	if directoryExists(configFilePath) {
		return fmt.Errorf("Config file already exists in the directory.")
	}

	// Create config file
	file, err := os.Create(configFilePath)
	if err != nil {
		return fmt.Errorf("Failed to create config file: %v", err)
	}
	defer file.Close()

	// Extract directory name
	_, dirName := filepath.Split(dirPath)

	// Create default configuration
	defaultConfig := Config{
		Name:  dirName,
		Email: email, // set email from parameter
	}

	// Encode default configuration as JSON
	configData, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("Failed to encode default configuration: %v", err)
	}

	// Write config data to file
	_, err = file.Write(configData)
	if err != nil {
		return fmt.Errorf("Failed to write default configuration to config file: %v", err)
	}

	fmt.Printf("Config file created: %s\n", configFilePath)

	return nil
}

func ValidateKeyPair(pubKey ed25519.PublicKey, privKey ed25519.PrivateKey) bool {
	// Create a new message
	message := []byte("test message")

	// Sign the message with the private key
	signature := ed25519.Sign(privKey, message)

	// Verify the signature with the public key
	verified := ed25519.Verify(pubKey, message, signature)

	if !verified {
		log.Printf("Key pair validation failed.")
	}

	return verified
}

func GenerateKeyPair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate key pair: %v", err)
	}

	// Validate the generated key pair
	if !ValidateKeyPair(pubKey, privKey) {
		return nil, nil, fmt.Errorf("Invalid key pair generated.")
	}

	return pubKey, privKey, nil
}

func WriteKeysToFile(pubKeyPath string, privKeyPath string, pubKey ed25519.PublicKey, privKey ed25519.PrivateKey) error {
	sshPublicKey, err := ssh.NewPublicKey(pubKey)
	if err != nil {
		return fmt.Errorf("Failed to create new SSH public key: %v", err)
	}

	sshPubKeyBytes := ssh.MarshalAuthorizedKey(sshPublicKey)

	pemBlock := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: edkey.MarshalED25519PrivateKey(privKey),
	}
	privKeyBytes := pem.EncodeToMemory(pemBlock)

	err = ioutil.WriteFile(privKeyPath, privKeyBytes, 0600)
	if err != nil {
		return fmt.Errorf("Failed to write private key to file: %v", err)
	}

	err = ioutil.WriteFile(pubKeyPath, sshPubKeyBytes, 0644)
	if err != nil {
		return fmt.Errorf("Failed to write public key to file: %v", err)
	}

	return nil
}

func directoryExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	log.Printf("Failed to check directory existence: %v", err)
	return false
}

// Config struct to hold the configuration data
type Config struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
