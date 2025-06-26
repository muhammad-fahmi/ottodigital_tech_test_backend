package database

import (
	"fmt"
	"log"
	"os"

	"my-backend-app/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() {
	var err error

	// Check if we're in test mode
	if os.Getenv("TEST_MODE") == "true" {
		// Use SQLite for testing
		DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			log.Fatal("Failed to connect to test database:", err)
		}
	} else {
		// Use MySQL for production
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)

		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(
		&models.Brand{},
		&models.Voucher{},
		&models.Customer{},
		&models.Transaction{},
		&models.TransactionItem{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	if os.Getenv("TEST_MODE") != "true" {
		log.Println("Database connected and migrated successfully")
	}
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
