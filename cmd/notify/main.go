package main

import (
	"github.com/JaKu01/GoNotify/internal"
	"github.com/JaKu01/GoNotify/web"
	"log"
)

func main() {

	log.Println("Starting GoNotify")
	log.Println("Initializing database")
	err := internal.InitDatabase()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	log.Println("Setting VAPID keys")
	err = internal.SetVapidKeys()
	if err != nil {
		log.Fatalf("Error setting VAPID keys: %v", err)
	}

	err = web.StartWebService()
	if err != nil {
		log.Fatalf("Error starting web service: %v", err)
	}
}
