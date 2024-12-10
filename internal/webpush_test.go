package internal

import (
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

var (
	testDB *gorm.DB
)

func setupTestDB(t *testing.T) {
	// Set up the in-memory SQLite database for testing
	tempDir := t.TempDir()
	dbPath := tempDir + "/test.db"
	var err error
	testDB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	err = testDB.AutoMigrate(&WebPushSubscription{})
	if err != nil {
		t.Fatalf("Failed to migrate database schema: %v", err)
	}

	// Set the global Connection variable to the in-memory DB
	Connection = testDB
}

// TestSaveSubscription tests the SaveSubscription function
func TestSaveSubscription(t *testing.T) {
	setupTestDB(t)

	subscriptionStr := []byte(`{"endpoint": "https://example.com", "keys": {"p256dh": "key1", "auth": "key2"}}`)

	// Test successful save
	err := SaveSubscription(subscriptionStr)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	// Check if the subscription exists in the database
	var dbSubscription WebPushSubscription
	err = testDB.Where("endpoint = ?", "https://example.com").First(&dbSubscription).Error
	if err != nil {
		t.Fatalf("Expected subscription to be saved, but got error: %v", err)
	}

	// Ensure that the saved subscription matches the input
	if dbSubscription.Endpoint != "https://example.com" {
		t.Fatalf("Expected endpoint to be 'https://example.com', but got: %v", dbSubscription.Endpoint)
	}

	// Test with invalid subscriptionStr (could be invalid JSON)
	invalidSubscriptionStr := []byte(`invalid-json`)
	err = SaveSubscription(invalidSubscriptionStr)
	if err != nil {
		t.Fatalf("Expected no error when invalid subscription is passed, but got: %v", err)
	}
}

// TestRemoveSubscription tests the RemoveSubscription function
func TestRemoveSubscription(t *testing.T) {
	setupTestDB(t)

	unsubscriptionRequest := WebPushUnsubscriptionRequest{
		Endpoint: "https://example.com",
	}

	// First, save a subscription to the database
	subscriptionStr := []byte(`{"endpoint": "https://example.com", "keys": {"p256dh": "key1", "auth": "key2"}}`)
	subscriptionForDb := WebPushSubscription{
		Endpoint:     unsubscriptionRequest.Endpoint,
		Subscription: string(subscriptionStr),
	}

	result := testDB.Create(&subscriptionForDb)
	if result.Error != nil {
		t.Fatalf("Failed to save subscription to the database: %v", result.Error)
	}

	// Test successful removal
	err := RemoveSubscription(unsubscriptionRequest)
	if err != nil {
		t.Fatalf("Expected no error during removal, but got: %v", err)
	}

	// Check if the subscription was removed from the database
	var dbSubscription WebPushSubscription
	err = testDB.Where("endpoint = ?", "https://example.com").First(&dbSubscription).Error
	if err == nil {
		t.Fatalf("Expected subscription to be removed, but it still exists")
	}
}
