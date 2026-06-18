package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"rentflow-webhook/internal/firebase"
	"rentflow-webhook/internal/models"
	"rentflow-webhook/internal/services"
)

// VerifyWebhook handles the webhook verification challenge from Meta
func VerifyWebhook(c *gin.Context, log *logrus.Logger, verifyToken string) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	log.WithFields(logrus.Fields{
		"mode":      mode,
		"token":     token,
		"challenge": challenge,
	}).Info("VerifyWebhook called")

	if mode == "subscribe" && token == verifyToken {
		log.Info("Webhook verified successfully")
		c.String(http.StatusOK, challenge)
	} else {
		log.Warn("Webhook verification failed")
		c.AbortWithStatus(http.StatusForbidden)
	}
}

// ProcessWebhook handles incoming webhook notifications from Meta
func ProcessWebhook(c *gin.Context, log *logrus.Logger, whatsappService *services.WhatsappService, firestoreService *firebase.FirestoreService) {
	log.Info("ProcessWebhook called")
	var webhook models.WhatsappWebhook
	if err := c.ShouldBindJSON(&webhook); err != nil {
		log.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// We only process the first change and first message for simplicity
	if len(webhook.Entry) > 0 && len(webhook.Entry[0].Changes) > 0 && len(webhook.Entry[0].Changes[0].Value.Messages) > 0 {
		message := webhook.Entry[0].Changes[0].Value.Messages[0]
		log = log.WithFields(logrus.Fields{
			"messageID": message.ID,
			"from":      message.From,
			"type":      message.Type,
		})

		if message.Type == "interactive" && message.Interactive.Type == "button_reply" {
			log.Info("Routing to HandleInteractiveMessage")
			HandleInteractiveMessage(log, whatsappService, firestoreService, &message)
		} else if message.Type == "text" {
			log.Info("Routing to HandleTextMessage")
			HandleTextMessage(c, log, whatsappService, firestoreService, &webhook)
		} else {
			log.Info("Ignoring non-interactive or non-text message")
		}

	} else {
		log.Info("Received webhook with no relevant message changes")
	}

	log.Info("Webhook processed successfully")
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
