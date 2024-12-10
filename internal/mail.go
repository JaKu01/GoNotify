package internal

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
)

var (
	smtpHost      = os.Getenv("SMTP_HOST")
	smtpPort      = os.Getenv("SMTP_PORT")
	emailAddress  = os.Getenv("EMAIL")
	emailPassword = os.Getenv("EMAIL_PASSWORD")
	tlsConfig     = &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpHost,
	}
)

// SendMail sends an email using SSL/TLS
func SendMail(request NotificationRequest) error {
	subject, contentType, body := extractEmailDetails(request)

	// Dial SMTP server
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// Create new SMTP client from the connection
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// Authenticate with the SMTP server
	auth := smtp.PlainAuth("", emailAddress, emailPassword, smtpHost)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Set the sender and recipient
	if err = client.Mail(emailAddress); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err = client.Rcpt(emailAddress); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send the email data
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	_, err = w.Write(createMessage(emailAddress, subject, contentType, body))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	// Quit the SMTP session
	if err = client.Quit(); err != nil {
		return fmt.Errorf("failed to close SMTP client: %w", err)
	}

	return nil
}

func extractEmailDetails(request NotificationRequest) (subject string, contentType string, body string) {
	if len(contentType) == 0 {
		contentType = "text/plain"
	}

	return request.Subject, contentType, request.Body
}

// createMessage constructs the email message
func createMessage(emailAddress, subject, contentType, body string) []byte {
	return []byte("From: " + emailAddress + "\r\n" +
		"To: " + emailAddress + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: " + contentType + "; charset=UTF-8\r\n\r\n" + // Add the Content-Type header for HTML
		body + "\r\n")
}
