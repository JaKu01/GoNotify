package web

import (
	"encoding/json"
	"github.com/JaKu01/GoNotify/internal"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		generateAndSendResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		// Serve static file
		handleStatic(w, r)
		return
	}

	// Parse and execute the template
	tmpl, err := template.ParseFiles("./template/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Template parsing error: %v", err)
		return
	}

	// Execute the template with the single string
	err = tmpl.Execute(w, internal.VapidPublicKey)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
	}
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	// Serve static file
	filePath := "./static" + r.URL.Path
	if _, err := os.Stat(filePath); err == nil {
		http.ServeFile(w, r, filePath)
		return
	} else {
		// File not found, return 404
		generateAndSendResponse(w, "The requested static file was not found.", http.StatusNotFound)
		return
	}
}

func handleSubscribe(w http.ResponseWriter, r *http.Request) {
	// extract body from request as string
	subscription, err := io.ReadAll(r.Body)
	if err != nil {
		generateAndSendResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	internal.SaveSubscription(string(subscription))
	generateAndSendResponse(w, "Subscription saved successfully", http.StatusOK)
}

func handleDeleteSubscribe(w http.ResponseWriter, r *http.Request) {
	// extract body from request as string
	var req internal.WebPushUnsubscriptionRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		generateAndSendResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = internal.RemoveSubscription(req)

	if err != nil {
		generateAndSendResponse(w, "An error occurred while deleting the subscription", http.StatusInternalServerError)
		return
	}

	generateAndSendResponse(w, "Subscription saved successfully", http.StatusOK)
}

func handleMail(w http.ResponseWriter, r *http.Request) {
	// read the request body
	var req internal.NotificationRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		generateAndSendResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// send the email
	err = internal.SendMail(req)
	if err != nil {
		generateAndSendResponse(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	generateAndSendResponse(w, "Email sent successfully", http.StatusOK)
}

func handleWebPush(w http.ResponseWriter, r *http.Request) {
	// read the request body
	var req internal.NotificationRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		generateAndSendResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	internal.SendToAllSubscribers(req)
	generateAndSendResponse(w, "Web push notifications sent", http.StatusOK)
}

func handleAll(w http.ResponseWriter, r *http.Request) {
	generateAndSendResponse(w, "Not yet implemented", http.StatusNotImplemented)
}

func generateAndSendResponse(w http.ResponseWriter, message string, statusCode int) {
	response := internal.NotificationResponse{Message: message}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Printf("Error: failed to write response %v", err)
	}
}

func handleServiceWorker(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "./static/service-worker.js")
}

func handleManifest(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "./static/manifest.json")
}
