package internal

import "testing"

// TestExtractEmailDetails tests correct extraction of request details
func TestExtractEmailDetails(t *testing.T) {
	req := NotificationRequest{
		Subject:     "Test Subject",
		ContentType: "text/plain",
		Body:        "Test Body",
	}

	subject, contentType, body := extractEmailDetails(req)

	if subject != req.Subject {
		t.Fatalf("Expected subject to be %s, but got %s", req.Subject, subject)
	}

	if contentType != req.ContentType {
		t.Fatalf("Expected contentType to be %s, but got %s", req.ContentType, contentType)
	}

	if body != req.Body {
		t.Fatalf("Expected body to be %s, but got %s", req.Body, body)
	}
}

// TestCreateMessage_Plain tests createMessage with plain text content
func TestCreateMessage_Plain(t *testing.T) {
	emailAddress := "test@example.com"
	subject := "Test Subject"
	contentType := "text/plain"
	body := "This is the body of the email."

	expectedMessage := "From: " + emailAddress + "\r\n" +
		"To: " + emailAddress + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: " + contentType + "; charset=UTF-8\r\n\r\n" +
		body + "\r\n"

	message := createMessage(emailAddress, subject, contentType, body)

	// Check if the constructed message matches the expected message
	if string(message) != expectedMessage {
		t.Fatalf("Expected message to be:\n%s\nBut got:\n%s", expectedMessage, message)
	}
}

// TestCreateMessage_HTMLContent tests createMessage with HTML content
func TestCreateMessage_HTMLContent(t *testing.T) {
	emailAddress := "test@example.com"
	subject := "Test HTML Email"
	contentType := "text/html"
	body := "<html><body><h1>Hello, World!</h1></body></html>"

	expectedMessage := "From: " + emailAddress + "\r\n" +
		"To: " + emailAddress + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: " + contentType + "; charset=UTF-8\r\n\r\n" +
		body + "\r\n"

	message := createMessage(emailAddress, subject, contentType, body)

	// Check if the constructed message matches the expected message
	if string(message) != expectedMessage {
		t.Fatalf("Expected message to be:\n%s\nBut got:\n%s", expectedMessage, message)
	}
}
