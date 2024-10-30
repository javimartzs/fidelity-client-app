package database

import (
	"fidelity-client-app/config"
	"fidelity-client-app/models"

	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {

	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s",
		config.Vars.DBUser,
		config.Vars.DBPass,
		config.Vars.DBName,
		config.Vars.DBHost,
		config.Vars.DBPort,
	)

	var err error
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	DB.AutoMigrate(
		&models.User{},
		&models.Promotion{},
	)

	return DB
}
