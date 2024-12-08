package web

import (
	"log"
	"net/http"
)

func StartWebService() {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/subscribe", handleSubscribe)
	mux.HandleFunc("POST /api/mail", handleMail)
	mux.HandleFunc("POST /api/webpush", handleWebPush)
	mux.HandleFunc("POST /api/all", handleAll)
	mux.HandleFunc("GET /service-worker.js", handleServiceWorker)
	mux.HandleFunc("GET /", handleIndex)

	loggedMux := LoggingMiddleware(mux) // enable logging middleware

	log.Printf("Server running")
	log.Fatal(http.ListenAndServe(":8080", loggedMux))
}
