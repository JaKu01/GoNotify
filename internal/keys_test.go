package internal

import (
	"os"
	"testing"
)

// TestSetVapidKeys_GenerateNewKeys tests SetVapidKeys when no keys exist on disk
func TestSetVapidKeys_GenerateNewKeys(t *testing.T) {
	// Temporary directory for test keys
	tempDir := t.TempDir()
	keysDirectory = tempDir
	privateKeyPath = tempDir + "/private.key"
	publicKeyPath = tempDir + "/public.key"

	err := SetVapidKeys()
	if err != nil {
		t.Fatalf("SetVapidKeys returned an error: %v", err)
	}

	// Check if VAPID keys are non-empty
	if VapidPublicKey == "" {
		t.Fatalf("Expected VapidPublicKey to be non-empty, but it was empty")
	}

	if vapidPrivateKey == "" {
		t.Fatalf("Expected vapidPrivateKey to be non-empty, but it was empty")
	}
}

func TestSetVapidKeys_LoadExistingKeys(t *testing.T) {
	// Temporary directory for test keys
	tempDir := t.TempDir()
	privateKeyPath = tempDir + "/private.key"
	publicKeyPath = tempDir + "/public.key"

	// Write fake keys to disk
	err := os.WriteFile(privateKeyPath, []byte("fake-private-key"), 0600)
	if err != nil {
		t.Fatalf("Failed to write private key to temp file: %v", err)
	}

	err = os.WriteFile(publicKeyPath, []byte("fake-public-key"), 0600)
	if err != nil {
		t.Fatalf("Failed to write public key to temp file: %v", err)
	}

	err = SetVapidKeys()
	if err != nil {
		t.Fatalf("SetVapidKeys returned an error: %v", err)
	}

	// Check if VAPID keys contain the loaded key values
	if vapidPrivateKey != "fake-private-key" {
		t.Fatalf("Expected vapidPrivateKey to be non-empty, but it was empty")
	}
	if VapidPublicKey != "fake-public-key" {
		t.Fatalf("Expected VapidPublicKey to be non-empty, but it was empty")
	}

}

// TestLoadKeys_Success tests loadKeys when keys exist on disk
func TestLoadKeys_Success(t *testing.T) {
	// Temporary directory for test keys
	tempDir := t.TempDir()
	privateKeyPath = tempDir + "/private.key"
	publicKeyPath = tempDir + "/public.key"

	// Write fake keys to disk
	err := os.WriteFile(privateKeyPath, []byte("fake-private-key"), 0600)
	if err != nil {
		t.Fatalf("Failed to write private key to temp file: %v", err)
	}

	err = os.WriteFile(publicKeyPath, []byte("fake-public-key"), 0600)
	if err != nil {
		t.Fatalf("Failed to write public key to temp file: %v", err)
	}

	// Call loadKeys
	privateKey, publicKey, err := loadKeys()
	if err != nil {
		t.Fatalf("loadKeys returned an error: %v", err)
	}

	// Assert that keys are loaded
	if privateKey != "fake-private-key" {
		t.Fatalf("Expected vapidPrivateKey to be 'fake-private-key', but got: %s", vapidPrivateKey)
	}
	if publicKey != "fake-public-key" {
		t.Fatalf("Expected VapidPublicKey to be 'fake-public-key', but got: %s", VapidPublicKey)
	}

}

// TestGenerateKeys tests the generation of VAPID keys
func TestGenerateKeys(t *testing.T) {
	// Set up a temporary directory for test keys
	tempDir := t.TempDir()
	keysDirectory = tempDir
	privateKeyPath = tempDir + "/private.key"
	publicKeyPath = tempDir + "/public.key"

	// Call generateKeys
	privateKey, publicKey, err := generateKeys()
	if err != nil {
		t.Fatalf("generateKeys returned an error: %v", err)
	}

	// Check if files were created
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		t.Fatalf("Expected private key file to be created, but it was not found")
	}

	if _, err := os.Stat(publicKeyPath); os.IsNotExist(err) {
		t.Fatalf("Expected public key file to be created, but it was not found")
	}

	// Check if VAPID keys are non-empty
	if privateKey == "" {
		t.Fatalf("Expected privateKey to be non-empty, but it was empty")
	}

	if publicKey == "" {
		t.Fatalf("Expected publicKey to be non-empty, but it was empty")
	}
}
