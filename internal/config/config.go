package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application's configuration
type Config struct {
	WhatsappVerifyToken     string
	WhatsappAccessToken     string
	WhatsappPhoneNumberID   string
	ProjectID               string
	FirebaseCredentialsJSON string
	Port                    string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() (*Config, error) {
	// This will load the .env file if it exists, and do nothing if it doesn't.
	// This is perfect for supporting both local development and production environments.
	godotenv.Load()

	return &Config{
		WhatsappVerifyToken:     os.Getenv("WHATSAPP_VERIFY_TOKEN"),
		WhatsappAccessToken:     os.Getenv("WHATSAPP_ACCESS_TOKEN"),
		WhatsappPhoneNumberID:   os.Getenv("WHATSAPP_PHONE_NUMBER_ID"),
		ProjectID:               os.Getenv("PROJECT_ID"),
		FirebaseCredentialsJSON: os.Getenv("FIREBASE_CREDENTIALS_JSON"),
		Port:                    os.Getenv("PORT"),
	}, nil
}
