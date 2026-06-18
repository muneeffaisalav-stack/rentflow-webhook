package main

import (
	"github.com/gin-gonic/gin"
	"rentflow-webhook/internal/config"
	"rentflow-webhook/internal/firebase"
	"rentflow-webhook/internal/handlers"
	"rentflow-webhook/internal/logger"
	"rentflow-webhook/internal/services"
)

func main() {
	// Create a new logger
	log := logger.NewLogger()

	log.Info("Starting server...")

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create a new WhatsApp service
	whatsappService := services.NewWhatsappService(cfg.WhatsappAccessToken, cfg.WhatsappPhoneNumberID, log)

	// Create a new Firestore service
	firestoreService, err := firebase.NewFirestoreService(cfg, log)
	if err != nil {
		log.Fatalf("Failed to create Firestore service: %v", err)
	}

	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router
	router := gin.New()

	// Add a recovery middleware to handle panics
	router.Use(gin.Recovery())

	// Add a logger middleware
	router.Use(gin.Logger())

	// Setup routes
	log.Info("Setting up routes...")
	router.GET("/health", handlers.HealthCheck(log, firestoreService))
	router.GET("/whatsapp", handlers.VerifyWebhook(cfg, log))
	router.POST("/whatsapp", handlers.ProcessWebhook(log, whatsappService, firestoreService))

	// Start the server
	log.Infof("Server started on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
