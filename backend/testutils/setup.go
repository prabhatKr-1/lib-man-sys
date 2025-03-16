package testutils

import (
	"log"
	"os"
	"testing"

	"github.com/prabhatKr-1/lib-man-sys/backend/config"
	"github.com/prabhatKr-1/lib-man-sys/backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB initializes an in-memory SQLite database for testing
func SetupTestDB() {
	var err error
	config.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to initialize test database: %v", err)
	}

	log.Println("✅ Test database initialized successfully!")

	// Migrate tables
	err = config.DB.AutoMigrate(
		&models.Library{},
		&models.User{},
		&models.Books{},
		&models.IssueRegistry{},
		&models.RequestEvents{},
	)
	if err != nil {
		log.Fatalf("❌ Failed to migrate test database: %v", err)
	}

	log.Println("✅ Test database migration successful")
}

// Run setup before executing tests
func TestMain(m *testing.M) {
	SetupTestDB()
	exitVal := m.Run()
	log.Println("✅ Tests completed")
	os.Exit(exitVal)
}
