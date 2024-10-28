package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Vars Config

type Config struct {
	DBUser string
	DBPass string
	DBName string
	DBHost string
	DBPort string
	JwtKey string
}

func LoadEnv() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading env file: %v", err)
	}

	Vars = Config{
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASS"),
		DBName: os.Getenv("DB_NAME"),
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		JwtKey: os.Getenv("JWT_KEY"),
	}

	fmt.Println("Environments var imported")

}
