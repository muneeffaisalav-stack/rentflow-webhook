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
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		WhatsappVerifyToken:     os.Getenv("WHATSAPP_VERIFY_TOKEN"),
		WhatsappAccessToken:     os.Getenv("WHATSAPP_ACCESS_TOKEN"),
		WhatsappPhoneNumberID:   os.Getenv("WHATSAPP_PHONE_NUMBER_ID"),
		ProjectID:               os.Getenv("PROJECT_ID"),
		FirebaseCredentialsJSON: os.Getenv("FIREBASE_CREDENTIALS_JSON"),
		Port:                    os.Getenv("PORT"),
	}, nil
}
