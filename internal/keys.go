package internal

import (
	"github.com/SherClockHolmes/webpush-go"
	"log"
	"os"
)

const (
	privateKeyPath = "./keys/private.key"
	publicKeyPath  = "./keys/public.key"
)

var (
	VapidPublicKey  string
	vapidPrivateKey string
)

// SetVapidKeys reads VAPID keys from disk or generates new ones if they don't exist
func SetVapidKeys() {
	if loadKeys() == nil {
		// loading keys was successful
		log.Println("VAPID keys loaded from disk")
		return
	}
	vapidPrivateKey, VapidPublicKey = generateKeys()
	log.Println("VAPID keys generated and saved to disk")
}

func loadKeys() error {
	publicKey, err := os.ReadFile(publicKeyPath)

	if err != nil {
		return err
	}

	privateKey, err := os.ReadFile(privateKeyPath)

	if err != nil {
		return err
	}

	VapidPublicKey = string(publicKey)
	vapidPrivateKey = string(privateKey)

	return nil
}

func generateKeys() (privateKey string, publicKey string) {
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		log.Fatalf("error generating VAPID keys: %v", err)
	}

	// Ensure the "keys" directory exists
	err = os.MkdirAll("./keys", os.ModePerm)
	if err != nil {
		log.Fatalf("error creating keys directory: %v", err)
	}

	// Write the private key to disk
	err = os.WriteFile(privateKeyPath, []byte(privateKey), 0600)
	if err != nil {
		log.Fatalf("error writing private key to file: %v", err)
	}

	// Write the public key to disk
	err = os.WriteFile(publicKeyPath, []byte(publicKey), 0600)
	if err != nil {
		log.Fatalf("error writing public key to file: %v", err)
	}

	return privateKey, publicKey
}
