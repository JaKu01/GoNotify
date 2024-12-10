package internal

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var (
	Connection   *gorm.DB
	DatabasePath = "./internal/sqlite/subscription.db"
)

func InitDatabase() error {
	var err error
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
