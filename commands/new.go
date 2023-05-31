package commands

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
	return "new <name>"
}

func (*NewCmd) SetFlags(f *flag.FlagSet) {}

func (*NewCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() < 1 {
		log.Printf("You must provide a name for the new directory.")
		return subcommands.ExitFailure
	}
	dirname := f.Arg(0)
	newDir := filepath.Join(config.GetKeysDir(), dirname)
	config.CreateDirIfNotExist(newDir)

	pubKeyPath := filepath.Join(newDir, "id_ed25519.pub")
	privKeyPath := filepath.Join(newDir, "id_ed25519")

	pubKey, privKey, err := GenerateKeyPair()
	if err != nil {
		log.Printf("%v", err)
		return subcommands.ExitFailure
	}

	err = WriteKeysToFile(pubKeyPath, privKeyPath, pubKey, privKey)
	if err != nil {
		log.Printf("%v", err)
		return subcommands.ExitFailure
	}

	fmt.Printf("Public key saved to: %s\n", pubKeyPath)
	fmt.Printf("Private key saved to: %s\n", privKeyPath)

	return subcommands.ExitSuccess
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
