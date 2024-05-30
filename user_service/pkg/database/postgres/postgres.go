package postgres

import (
	"fmt"
	"os"
	"time"

	"github.com/catness812/faf-hub-backend/user_service/internal/models"
	"github.com/gookit/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadDatabase() *gorm.DB {
	db := connect()
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		slog.Error(err)
	}
	return db
}

func connect() *gorm.DB {
	var err error

	dsn := fmt.Sprintf(`host=%s dbname=%s user=%s password=%s port=%s sslmode=disable`,
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
	)

	for i := 0; i < 5; i++ {
		database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			time.Sleep(3 * time.Second)
		} else {
			slog.Info("Successfully connected to the Postgres database")
			return database
		}
	}

	slog.Fatalf("Could not connect to the database after %d attempts: %v", 3, err)
	return nil
}
