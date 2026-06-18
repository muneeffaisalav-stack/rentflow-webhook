package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"rentflow-webhook/internal/firebase"
	"rentflow-webhook/internal/models"
	"rentflow-webhook/internal/services"
)

// HandleTextMessage handles incoming text messages
func HandleTextMessage(c *gin.Context, log *logrus.Logger, whatsappService *services.WhatsappService, firestoreService *firebase.FirestoreService, webhook *models.WhatsappWebhook) {
	log.Info("Handling text message")
	message := webhook.Entry[0].Changes[0].Value.Messages[0]
	from := message.From
	textBody := message.Text.Body

	log = log.WithFields(logrus.Fields{
		"from":     from,
		"textBody": textBody,
	})

	// Check if the message is 'pay'
	if strings.ToLower(textBody) == "pay" {
		log.Info("Received 'pay' command, initiating test payment flow")

		invoiceID := "mock_invoice_123"
		tenantName := "John Doe"
		landlordPhone := from
		amount := "1000"
		landlordName := "Jane Smith"

		text := fmt.Sprintf(
			"A payment of $%s for invoice %s has been received from %s. Please confirm you have received the payment.",
			amount, invoiceID, tenantName,
		)

		buttons := []models.Button{
			{
				Type: "reply",
				Reply: models.ButtonReply{
					ID:    fmt.Sprintf("mark_paid_%s_%s", invoiceID, landlordName),
					Title: "Mark as Paid",
				},
			},
		}

		log.Info("Sending test button message")
		err := whatsappService.SendButtonMessage(landlordPhone, text, buttons)
		if err != nil {
			log.WithError(err).Error("Failed to send test button message")
		} else {
			log.Info("Test payment notification sent successfully")
		}

	} else {
		log.Info("Ignoring text message that is not a test command")
	}
}
