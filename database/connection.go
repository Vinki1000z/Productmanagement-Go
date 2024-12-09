package database

import (
	"log"
	"productManagmentBackend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase connects to the PostgreSQL database and returns an error if it fails
func ConnectDatabase() error {
	dsn := "host=localhost user=postgres password=Vikrant@1234 dbname=products port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = database

	// Automatically migrate the schema
	err = DB.AutoMigrate(&models.User{}, &models.Product{})
	if err != nil {
		return err
	}

	log.Println("Database connected and tables migrated successfully")
	return nil
}
