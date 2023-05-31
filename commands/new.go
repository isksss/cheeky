package commands

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/subcommands"
	"github.com/isksss/cheeky/config"
	"golang.org/x/crypto/ed25519"
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

	pub, priv, err := generateKeyPair()
	if err != nil {
		log.Printf("Failed to generate key pair: %v", err)
		return subcommands.ExitFailure
	}

	pubKeyPath := filepath.Join(newDir, "pub.key")
	privKeyPath := filepath.Join(newDir, "priv.key")

	if err := saveKeyToFile(pubKeyPath, pub, 0644); err != nil {
		log.Printf("Failed to save public key: %v", err)
		return subcommands.ExitFailure
	}

	if err := saveKeyToFile(privKeyPath, priv, 0600); err != nil {
		log.Printf("Failed to save private key: %v", err)
		return subcommands.ExitFailure
	}

	fmt.Printf("Public key saved to: %s\n", pubKeyPath)
	fmt.Printf("Private key saved to: %s\n", privKeyPath)

	return subcommands.ExitSuccess
}

func generateKeyPair() ([]byte, []byte, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	pubKeyEncoded := encodeKey(pub)
	privKeyEncoded := encodeKey(priv)
	return pubKeyEncoded, privKeyEncoded, nil
}

func encodeKey(key []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(key))
}

func saveKeyToFile(path string, key []byte, perm os.FileMode) error {
	pubTemplate, privTemplate, err := config.GetTemplates()
	if err != nil {
		return err
	}

	type KeyData struct {
		Key string
	}
	data := KeyData{
		Key: string(key),
	}

	var buf bytes.Buffer
	if strings.Contains(path, "pub.key") {
		err = pubTemplate.Execute(&buf, data)
		if err != nil {
			return err
		}
	} else {
		err = privTemplate.Execute(&buf, data)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(path, buf.Bytes(), perm)
}
