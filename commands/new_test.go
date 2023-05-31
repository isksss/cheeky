package commands

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	pubKey, privKey, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}
	if pubKey == nil || privKey == nil {
		t.Fatal("Failed to generate key pair: public or private key is nil")
	}
	if !ValidateKeyPair(pubKey, privKey) {
		t.Fatal("Generated key pair is invalid.")
	}
}

func TestValidateKeyPair(t *testing.T) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}
	if !ValidateKeyPair(pubKey, privKey) {
		t.Fatal("Valid key pair is considered invalid.")
	}
}
