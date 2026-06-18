package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"rentflow-webhook/internal/config"
	"rentflow-webhook/internal/firebase"
	"rentflow-webhook/internal/models"
	"rentflow-webhook/internal/services"
)

// VerifyWebhook verifies the webhook
func VerifyWebhook(cfg *config.Config, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("VerifyWebhook called")
		mode := c.Query("hub.mode")
		token := c.Query("hub.verify_token")
		challenge := c.Query("hub.challenge")

		log.WithFields(logrus.Fields{
			"mode":      mode,
			"token":     token,
			"challenge": challenge,
		}).Info("Received verification request")

		if mode == "subscribe" && token == cfg.WhatsappVerifyToken {
			log.Info("Webhook verification successful")
			c.Writer.WriteString(challenge)
		} else {
			log.WithFields(logrus.Fields{
				"mode":          mode,
				"token":         token,
				"expectedToken": cfg.WhatsappVerifyToken,
			}).Warn("Webhook verification failed")
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
}

// ProcessWebhook processes incoming webhooks
func ProcessWebhook(log *logrus.Logger, whatsappService *services.WhatsappService, firestoreService *firebase.FirestoreService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("ProcessWebhook called")
		var webhook models.WhatsappWebhook
		if err := c.ShouldBindJSON(&webhook); err != nil {
			log.WithError(err).Error("Failed to parse webhook body")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if len(webhook.Entry) > 0 && len(webhook.Entry[0].Changes) > 0 && len(webhook.Entry[0].Changes[0].Value.Messages) > 0 {
			message := webhook.Entry[0].Changes[0].Value.Messages[0]

			// Check if it's an interactive message reply or a button click
			if (message.Interactive != nil && message.Interactive.Type == "button_reply") || message.Type == "button" {
				HandleInteractiveMessage(log, whatsappService, firestoreService, &message)
			} else {
				log.WithField("type", message.Type).Info("Ignoring non-interactive or non-button message")
			}
		} else {
			log.Info("Ignoring empty or non-message webhook")
		}

		log.Info("Webhook processed successfully")
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
