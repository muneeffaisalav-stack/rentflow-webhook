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
	var buttonReplyID string
	if message.Interactive != nil {
		buttonReplyID = message.Interactive.ButtonReply.ID
	} else if message.Button != nil {
		buttonReplyID = message.Button.Payload
	}

	senderPhoneNumber := message.From // Get the phone number of the person who clicked the button

	if strings.HasPrefix(buttonReplyID, "mark_paid_") {
		invoiceID := strings.TrimPrefix(buttonReplyID, "mark_paid_")
		log := log.WithField("invoiceID", invoiceID)

		// --- Security Check for Tenant --- //
		invoice, err := firestoreService.GetInvoice(context.Background(), invoiceID)
		if err != nil {
			log.WithError(err).Error("Failed to get invoice for security check")
			return
		}
		tenant, err := firestoreService.GetTenant(context.Background(), invoice.TenantID)
		if err != nil {
			log.WithError(err).Error("Failed to get tenant for security check")
			return
		}

		// VERIFY: Is the person clicking the button the actual tenant?
		if strings.TrimPrefix(senderPhoneNumber, "+") != strings.TrimPrefix(tenant.Phone, "+") {
			log.Warnf("Unauthorized 'Mark as Paid' attempt. Sender: %s, Expected Tenant: %s", senderPhoneNumber, tenant.Phone)
			return // Stop processing immediately
		}

		// --- End Security Check --- //

		log.Info("'Mark as Paid' button clicked by authorized tenant")

		landlord, err := firestoreService.GetUser(context.Background(), invoice.LandlordID)
		if err != nil {
			log.WithError(err).Error("Failed to get landlord")
			return
		}

		property, err := firestoreService.GetProperty(context.Background(), invoice.PropertyID)
		if err != nil {
			log.WithError(err).Error("Failed to get property")
			return
		}

		// --- Send Landlord Verification Template ---
		templateName := "payment_verification_v1"
		languageCode := "en_US" // Or your default language code

		// Construct the components for the template message
		components := []map[string]interface{}{
			{
				"type": "body",
				"parameters": []map[string]interface{}{
					{"type": "text", "text": landlord.Name},
					{"type": "text", "text": tenant.Name},
					{"type": "text", "text": fmt.Sprintf("%d", invoice.Amount)},
					{"type": "text", "text": invoice.Month},
					{"type": "text", "text": property.PropertyName},
				},
			},
			{
				"type": "button",
				"sub_type": "quick_reply",
				"index": "0",
				"parameters": []map[string]interface{}{
					{"type": "payload", "payload": fmt.Sprintf("received_%s", invoiceID)},
				},
			},
			{
				"type": "button",
				"sub_type": "quick_reply",
				"index": "1",
				"parameters": []map[string]interface{}{
					{"type": "payload", "payload": fmt.Sprintf("not_paid_%s", invoiceID)},
				},
			},
		}

		if err := whatsappService.SendTemplateMessage(landlord.PhoneNumber, templateName, languageCode, components); err != nil {
			log.WithError(err).Error("Failed to send verification template to landlord")
		}

	} else if strings.HasPrefix(buttonReplyID, "received_") || strings.HasPrefix(buttonReplyID, "not_paid_") {
		var action, status string
		var invoiceID string

		if strings.HasPrefix(buttonReplyID, "received_") {
			action = "Received"
			status = "paid"
			invoiceID = strings.TrimPrefix(buttonReplyID, "received_")
		} else {
			action = "Not Paid"
			status = "pending"
			invoiceID = strings.TrimPrefix(buttonReplyID, "not_paid_")
		}

		log := log.WithField("invoiceID", invoiceID)

		// --- Security Check for Landlord --- //
		invoice, err := firestoreService.GetInvoice(context.Background(), invoiceID)
		if err != nil {
			log.WithError(err).Error("Failed to get invoice for security check")
			return
		}
		landlord, err := firestoreService.GetUser(context.Background(), invoice.LandlordID)
		if err != nil {
			log.WithError(err).Error("Failed to get landlord for security check")
			return
		}

		// VERIFY: Is the person clicking the button the actual landlord?
		if strings.TrimPrefix(senderPhoneNumber, "+") != strings.TrimPrefix(landlord.PhoneNumber, "+") {
			log.Warnf("Unauthorized '%s' attempt. Sender: %s, Expected Landlord: %s", action, senderPhoneNumber, landlord.PhoneNumber)
			return // Stop processing immediately
		}
		// --- End Security Check --- //

		log.Infof("'%s' button clicked by authorized landlord", action)

		if err := firestoreService.UpdateInvoiceStatus(context.Background(), invoiceID, status); err != nil {
			log.WithError(err).Errorf("Failed to update invoice status to '%s'", status)
			return
		}
	}
}
