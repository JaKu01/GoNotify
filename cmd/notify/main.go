package main

import (
	web "github.com/JaKu01/GoNotify"
	"github.com/JaKu01/GoNotify/internal"
	"log"
)

func main() {

	log.Println("Starting GoNotify")
	log.Println("Initializing database")
	internal.InitDatabase()

	log.Println("Setting VAPID keys")
	internal.SetVapidKeys()

	web.StartWebService()
}
