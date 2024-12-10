package internal

import (
	"fmt"
	"github.com/SherClockHolmes/webpush-go"
	"log"
	"os"
)

var (
	keysDirectory   = "./keys"
	privateKeyPath  = keysDirectory + "/private.key"
	publicKeyPath   = keysDirectory + "/public.key"
	vapidPrivateKey string
	VapidPublicKey  string
)

// SetVapidKeys reads VAPID keys from disk or generates new ones if they don't exist
func SetVapidKeys() error {
	privateKey, publicKey, err := loadKeys()
	if err != nil {
		// loading keys was successful
		privateKey, publicKey, err = generateKeys()
		if err != nil {
			return fmt.Errorf("error generating VAPID keys: %v", err)
		}
	}

	VapidPublicKey = publicKey
	vapidPrivateKey = privateKey

	log.Println("VAPID keys generated and saved to disk")
	return nil
}

func loadKeys() (privateKey string, publicKey string, err error) {
	privateKeyFromFile, err := os.ReadFile(privateKeyPath)

	if err != nil {
		return privateKey, publicKey, err
	}

	publicKeyFromFile, err := os.ReadFile(publicKeyPath)

	if err != nil {
		return privateKey, publicKey, err
	}

	privateKey = string(privateKeyFromFile)
	publicKey = string(publicKeyFromFile)
	return privateKey, publicKey, nil
}

func generateKeys() (privateKey string, publicKey string, err error) {
	privateKey, publicKey, err = webpush.GenerateVAPIDKeys()
	if err != nil {
		return privateKey, publicKey, fmt.Errorf("error generating VAPID keys: %v", err)
	}

	// Ensure the "keys" directory exists
	err = os.MkdirAll(keysDirectory, os.ModePerm)
	if err != nil {
		return privateKey, publicKey, fmt.Errorf("error creating keys directory: %v", err)
	}

	// Write the private key to disk
	err = os.WriteFile(privateKeyPath, []byte(privateKey), 0600)
	if err != nil {
		return privateKey, publicKey, fmt.Errorf("error writing private key to file: %v", err)
	}

	// Write the public key to disk
	err = os.WriteFile(publicKeyPath, []byte(publicKey), 0600)
	if err != nil {
		return privateKey, publicKey, fmt.Errorf("error writing public key to file: %v", err)
	}

	return privateKey, publicKey, nil
}
