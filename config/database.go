package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"example.com/go-project/helper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost" // adjust as needed
	port     = 5432        // adjust as needed
	user     = "root"      // adjust as needed
	password = "root"      // adjust as needed
	dbname   = "test"      // adjust as needed
)

func DatabaseConnection() *gorm.DB {
	// Create the PostgreSQL connection string (DSN)
	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Set up custom logger for GORM
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Output to Stdout
		logger.Config{
			SlowThreshold:             time.Second, // Log queries that are slower than 1 second
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore 'record not found' errors
			Colorful:                  true,        // Enable colorful logs
		},
	)

	// Connect to the database
	db, err := gorm.Open(postgres.Open(sqlInfo), &gorm.Config{
		Logger: newLogger, // Attach the custom logger
	})
	helper.ErrorPanic(err)

	// Optional: Setup connection pooling (highly recommended for production)
	sqlDB, err := db.DB()
	helper.ErrorPanic(err)

	// Set max open connections to the database
	sqlDB.SetMaxOpenConns(100)

	// Set max idle connections
	sqlDB.SetMaxIdleConns(10)

	// Set max lifetime of a connection
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}
