package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	JwtSecret   string
}

func LoadConfig() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		log.Println("Failed to load env")
	}

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JwtSecret:   os.Getenv("JWT_SECRET"),
	}

	return cfg, nil

}
