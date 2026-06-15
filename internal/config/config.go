package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost string
	DBPort string
	DBName string
	DBUser string
	DBPass string

	// SMTP
	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string

	// Email
	EmailFrom      string
	EmailRecipients string

	// Server
	ServerPort string
	ServerEnv  string
}

func LoadConfig() (*Config, error) {
	godotenv.Load()

	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		smtpPort = 587
	}

	cfg := &Config{
		DBHost:          os.Getenv("DB_HOST"),
		DBPort:          os.Getenv("DB_PORT"),
		DBName:          os.Getenv("DB_NAME"),
		DBUser:          os.Getenv("DB_USER"),
		DBPass:          os.Getenv("DB_PASS"),
		SMTPHost:        os.Getenv("SMTP_HOST"),
		SMTPPort:        smtpPort,
		SMTPUser:        os.Getenv("SMTP_USER"),
		SMTPPass:        os.Getenv("SMTP_PASS"),
		EmailFrom:       os.Getenv("EMAIL_FROM"),
		EmailRecipients: os.Getenv("EMAIL_RECIPIENTS"),
		ServerPort:      os.Getenv("SERVER_PORT"),
		ServerEnv:       os.Getenv("SERVER_ENV"),
	}

	return cfg, nil
}
