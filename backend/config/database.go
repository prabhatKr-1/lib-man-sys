package config

import (
	"log"
	"os"

	"github.com/prabhatKr-1/lib-man-sys/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {

	// Postgres or NeonDB connection string
	dsn := os.Getenv("DATABASE_URL")

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db

	// Get the underlying SQL DB object
	sqlDB, err := db.DB()

	if err != nil {
		log.Fatalf("Failed to get DB object: %v", err)
	}

	// Verify connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	// Database migration
	err = DB.AutoMigrate(
		&models.Library{},
		&models.User{},
		&models.Books{},
		&models.IssueRegistry{},
		&models.RequestEvents{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Successfully connected to Postgres database!")

}
