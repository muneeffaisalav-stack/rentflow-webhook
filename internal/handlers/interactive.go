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
	buttonReplyID := message.Interactive.ButtonReply.ID

	if strings.HasPrefix(buttonReplyID, "mark_paid_") {
		invoiceID := strings.TrimPrefix(buttonReplyID, "mark_paid_")
		log.WithField("invoiceID", invoiceID).Info("Mark as Paid button clicked")

		// Fetch the invoice from Firestore
		invoice, err := firestoreService.GetInvoice(context.Background(), invoiceID)
		if err != nil {
			log.WithError(err).Error("Failed to get invoice")
			return
		}

		// Fetch tenant and landlord data
		tenant, err := firestoreService.GetUser(context.Background(), invoice.TenantRef.ID)
		if err != nil {
			log.WithError(err).Error("Failed to get tenant")
			return
		}
		landlord, err := firestoreService.GetUser(context.Background(), invoice.LandlordRef.ID)
		if err != nil {
			log.WithError(err).Error("Failed to get landlord")
			return
		}

		// Send a message to the landlord
		messageText := fmt.Sprintf("Tenant %s has marked the rent of %.2f as paid. Please verify.", tenant.Name, invoice.Amount)
		buttons := []models.Button{
			{Type: "reply", Title: "Received", ID: fmt.Sprintf("received_%s", invoiceID)},
			{Type: "reply", Title: "Not Paid", ID: fmt.Sprintf("not_paid_%s", invoiceID)},
		}
		if err := whatsappService.SendButtonMessage(landlord.PhoneNumber, messageText, buttons); err != nil {
			log.WithError(err).Error("Failed to send verification message to landlord")
		}

	} else if strings.HasPrefix(buttonReplyID, "received_") {
		invoiceID := strings.TrimPrefix(buttonReplyID, "received_")
		log.WithField("invoiceID", invoiceID).Info("Received button clicked")

		// Update the invoice status to "paid"
		if err := firestoreService.UpdateInvoiceStatus(context.Background(), invoiceID, "paid"); err != nil {
			log.WithError(err).Error("Failed to update invoice status")
			return
		}

	} else if strings.HasPrefix(buttonReplyID, "not_paid_") {
		invoiceID := strings.TrimPrefix(buttonReplyID, "not_paid_")
		log.WithField("invoiceID", invoiceID).Info("Not Paid button clicked")

		// Update the invoice status to "pending"
		if err := firestoreService.UpdateInvoiceStatus(context.Background(), invoiceID, "pending"); err != nil {
			log.WithError(err).Error("Failed to update invoice status")
			return
		}
	}
}
