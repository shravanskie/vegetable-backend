package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Provide safe defaults when environment variables are missing.
	dbUser := getEnvDefault("DB_USER", "root")
	dbPassword := getEnvDefault("DB_PASSWORD", "")
	dbHost := getEnvDefault("DB_HOST", "127.0.0.1")
	dbPort := getEnvDefault("DB_PORT", "3306")
	dbName := getEnvDefault("DB_NAME", "vegetable_app")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	var err error
	// Try connecting with retries/backoff to handle DB startup race conditions.
	maxAttempts := 10
	backoff := 500 // milliseconds
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Connected to MySQL successfully")
			return
		}
		log.Printf("attempt %d: failed to connect to DB: %v", attempt, err)
		if attempt < maxAttempts {
			sleepMs := backoff * (1 << (attempt - 1))
			if sleepMs > 30000 {
				sleepMs = 30000
			}
			// Sleep before next attempt
			time.Sleep(time.Duration(sleepMs) * time.Millisecond)
		}
	}
	log.Fatalf("failed to connect to DB after %d attempts: %v", maxAttempts, err)
}

// getEnvDefault returns the value of the environment variable named by the key.
// If the variable is empty or not present, it returns def.
func getEnvDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
