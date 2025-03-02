package internal

import (
	"fmt"
	"log"
	"os"
	"path"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	Connection   *gorm.DB
	DatabasePath = "./internal/sqlite/subscription.db"
)

func InitDatabase() error {
	// Ensure that all directories of the DatabasePath exists.
	// Otherwise, sqlite will fail when trying to open the database file, i.e. DatabasePath.
	databaseParentDirectoryPath := path.Dir(DatabasePath)
	err := os.MkdirAll(databaseParentDirectoryPath, 0775)
	if err != nil {
		return fmt.Errorf("failed to create parent directories '%s' of the database path: %w", databaseParentDirectoryPath, err)
	}

	Connection, err = gorm.Open(sqlite.Open(DatabasePath), &gorm.Config{})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	err = Connection.AutoMigrate(&WebPushSubscription{})

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database connection initialized")
	return nil
}
