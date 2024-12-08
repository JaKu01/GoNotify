package web

import (
	"encoding/json"
	"github.com/JaKu01/GoNotify/internal"
	"html/template"
	"io"
	"log"
	"net/http"
)

type MailNotificationRequest struct {
	Subject     string `json:"subject"`
	ContentType string `json:"content_type"`
	Body        string `json:"body"`
}

type NotificationResponse struct {
	Message string `json:"message"`
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

func handleMail(w http.ResponseWriter, r *http.Request) {
	// read the request body
	var req MailNotificationRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		generateAndSendResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// send the email
	err = internal.SendMail(req.Subject, req.Body)
	if err != nil {
		generateAndSendResponse(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	generateAndSendResponse(w, "Email sent successfully", http.StatusOK)
}

func handleWebPush(w http.ResponseWriter, r *http.Request) {
	// read the request body
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		generateAndSendResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	internal.SendToAllSubscribers(body)
	generateAndSendResponse(w, "Web push notifications sent", http.StatusOK)
}

func handleAll(w http.ResponseWriter, r *http.Request) {
	generateAndSendResponse(w, "Not yet implemented", http.StatusNotImplemented)
}

func generateAndSendResponse(w http.ResponseWriter, message string, statusCode int) {
	response := NotificationResponse{Message: message}

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
