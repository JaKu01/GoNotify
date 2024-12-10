package web

import (
	"github.com/JaKu01/GoNotify/internal"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleIndex(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	internal.VapidPublicKey = "test-public-key"

	// Create a response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleIndex)
	handler.ServeHTTP(rr, req)

	// Check if the status code is OK (200)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code 200, got %v", status)
	}

	// Check if the response contains the VapidPublicKey
	expected := internal.VapidPublicKey
	body := rr.Body.String()

	if !strings.Contains(body, expected) {
		t.Errorf("expected body to contain %v, got %v", expected, rr.Body.String())
	}
}

func TestHandleStatic(t *testing.T) {
	req, err := http.NewRequest("GET", "/subscribe.js", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleStatic)

	handler.ServeHTTP(rr, req)

	// Check for status code 200
	if rr.Code != http.StatusOK {
		t.Errorf("expected status code 404, got %v", rr.Code)
	}
}
