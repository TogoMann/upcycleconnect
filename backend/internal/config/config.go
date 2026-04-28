package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort  string
	DBHost   string
	DBPort   string
	DBUser   string
	DBPass   string
	DBName   string
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
	FromEmail string
}

func Load() *Config {
	godotenv.Load()
	return &Config{
		AppPort:   os.Getenv("APP_PORT"),
		DBHost:    os.Getenv("DB_HOST"),
		DBPort:    os.Getenv("DB_PORT"),
		DBUser:    os.Getenv("DB_USER"),
		DBPass:    os.Getenv("DB_PASSWORD"),
		DBName:    os.Getenv("DB_NAME"),
		SMTPHost:  os.Getenv("SMTP_HOST"),
		SMTPPort:  os.Getenv("SMTP_PORT"),
		SMTPUser:  os.Getenv("SMTP_USER"),
		SMTPPass:  os.Getenv("SMTP_PASS"),
		FromEmail: os.Getenv("FROM_EMAIL"),
	}
}
