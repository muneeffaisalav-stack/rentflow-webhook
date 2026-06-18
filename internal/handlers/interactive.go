package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"rentflow-webhook/internal/firebase"
	"rentflow-webhook/internal/models"
	"rentflow-webhook/internal/services"
)

func HandleInteractiveMessage(log *logrus.Logger, whatsappService *services.WhatsappService, firestoreService *firebase.FirestoreService, message *models.Message) {
	log.Info("Handling interactive message")
	buttonReplyID := message.Interactive.ButtonReply.ID
	senderPhoneNumber := message.From
	log = log.WithFields(logrus.Fields{
		"buttonReplyID":     buttonReplyID,
		"senderPhoneNumber": senderPhoneNumber,
	})

	if strings.HasPrefix(buttonReplyID, "mark_paid_") {
		log.Info("Processing 'mark_paid' action")
		invoiceID := strings.TrimPrefix(buttonReplyID, "mark_paid_")
		log = log.WithField("invoiceID", invoiceID)

		invoice, err := firestoreService.GetInvoice(context.Background(), invoiceID)
		if err != nil {
			return // Error is already logged in firestoreService
		}

		tenant, err := firestoreService.GetUser(context.Background(), invoice.TenantID)
		if err != nil {
			return // Error is already logged in firestoreService
		}

		if senderPhoneNumber != tenant.PhoneNumber {
			log.Warn("Unauthorized 'Mark as Paid' attempt")
			return
		}
		log.Info("Authorized tenant confirmed")

		landlord, err := firestoreService.GetUser(context.Background(), invoice.LandlordID)
		if err != nil {
			return // Error is already logged in firestoreService
		}

		templateName := "payment_verification_v1"
		log.WithFields(logrus.Fields{"templateName": templateName, "landlordPhone": landlord.PhoneNumber}).Info("Sending verification template to landlord")

		components := []map[string]interface{}{ /* ... Omitted for brevity ... */ }

		if err := whatsappService.SendTemplateMessage(landlord.PhoneNumber, templateName, "en_US", components); err != nil {
			log.WithError(err).Error("Failed to send verification template")
		} else {
			log.Info("Successfully sent verification template")
		}

	} else if strings.HasPrefix(buttonReplyID, "received_") || strings.HasPrefix(buttonReplyID, "not_paid_") {
		var action, status, invoiceID string
		if strings.HasPrefix(buttonReplyID, "received_") {
			action, status, invoiceID = "Received", "paid", strings.TrimPrefix(buttonReplyID, "received_")
		} else {
			action, status, invoiceID = "Not Paid", "pending", strings.TrimPrefix(buttonReplyID, "not_paid_")
		}
		log = log.WithFields(logrus.Fields{"action": action, "status": status, "invoiceID": invoiceID})
		log.Info("Processing landlord payment confirmation")

		invoice, err := firestoreService.GetInvoice(context.Background(), invoiceID)
		if err != nil {
			return
		}

		landlord, err := firestoreService.GetUser(context.Background(), invoice.LandlordID)
		if err != nil {
			return
		}

		if senderPhoneNumber != landlord.PhoneNumber {
			log.Warn("Unauthorized payment confirmation attempt")
			return
		}
		log.Info("Authorized landlord confirmed")

		if err := firestoreService.UpdateInvoiceStatus(context.Background(), invoiceID, status); err != nil {
			return // Error logged in service
		}
		log.Info("Successfully updated invoice status")
	} else {
		log.Warn("Unknown interactive button prefix")
	}
}
