package web

import (
	"log"
	"net/http"
)

func StartWebService() error {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/subscribe", handleSubscribe)
	mux.HandleFunc("DELETE /api/subscribe", handleDeleteSubscribe)
	mux.HandleFunc("POST /api/mail", handleMail)
	mux.HandleFunc("POST /api/webpush", handleWebPush)
	mux.HandleFunc("POST /api/telegram", handleTelegram)
	mux.HandleFunc("POST /api/all", handleAll)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/", handleIndex)

	loggedMux := LoggingMiddleware(mux) // enable logging middleware

	log.Printf("Server running")
	err := http.ListenAndServe(":8080", loggedMux)
	return err
}
