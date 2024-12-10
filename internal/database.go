package internal

import (
	"github.com/SherClockHolmes/webpush-go"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var Connection *gorm.DB

func InitDatabase() {
	var err error
	Connection, err = gorm.Open(sqlite.Open("./internal/sqlite/test.db"), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = Connection.AutoMigrate(&WebPushSubscription{})

	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("Database connection initialized")
}
