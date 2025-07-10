package db

import (
	calculationservice "GO-Calc/internal/calculationService"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() (*gorm.DB, error) {
	// TODO вынести в env
	// data source name
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5432 sslmode=disable" // sslmode безопасное соединение

	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}) // если менять в настройках базы данных

	if err != nil {
		log.Fatalf("Could not connect to databse: %v",  err)
	}

	// нет SQL . автомиграция
	if err := db.AutoMigrate(&calculationservice.Calculation{}); err != nil {
		log.Fatalf("Could not migrate database: %v", err)
	}

	return db, nil
}