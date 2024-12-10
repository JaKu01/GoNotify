package internal

import (
	"testing"
)

func TestInitDatabase(t *testing.T) {
	tempDir := t.TempDir()
	tempDBPath := tempDir + "/test.db"

	DatabasePath = tempDBPath

	err := InitDatabase()
	if err != nil {
		t.Fatalf("InitDatabase returned an error: %v", err)
	}

	if Connection == nil {
		t.Fatalf("Expected a non-nil database connection, but got nil")
	}

	sqlDB, err := Connection.DB()
	if err != nil {
		t.Fatalf("Failed to access the underlying SQL DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("Failed to ping the database: %v", err)
	}

	if !Connection.Migrator().HasTable(&WebPushSubscription{}) {
		t.Fatalf("Expected the WebPushSubscription table to exist, but it does not")
	}
}
