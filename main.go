package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"rentflow-webhook/internal/config"
	"rentflow-webhook/internal/firebase"
	"rentflow-webhook/internal/handlers"
	"rentflow-webhook/internal/services"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Firestore service
	firestoreService, err := firebase.NewFirestoreService(cfg, logger)
	if err != nil {
		logger.Fatalf("Failed to create Firestore service: %v", err)
	}

	// Initialize services
	whatsappService := services.NewWhatsappService(cfg.WhatsappAccessToken, cfg.WhatsappPhoneNumberID, logger)

	// Initialize Gin router
	router := gin.Default()

	// Set up routes
	router.GET("/webhook", handlers.VerifyWebhook(cfg, logger))
	router.POST("/webhook", handlers.ProcessWebhook(logger, whatsappService, firestoreService))

	// Start the server
	port := ":" + cfg.Port
	if err := router.Run(port); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
