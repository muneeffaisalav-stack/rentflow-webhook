package firebase

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"rentflow-webhook/internal/models"
)

// FirestoreService is a service for interacting with Firestore
type FirestoreService struct {
	client *firestore.Client
	log    *logrus.Logger
}

// NewFirestoreService creates a new FirestoreService
func NewFirestoreService(ctx context.Context, projectID string, serviceAccount []byte, log *logrus.Logger) (*FirestoreService, error) {
	conf := &firebase.Config{ProjectID: projectID}
	opt := option.WithCredentialsJSON(serviceAccount)

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.WithError(err).Error("Failed to create Firebase app")
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to create Firestore client")
		return nil, fmt.Errorf("error creating client: %v", err)
	}

	log.Info("Firestore service created successfully")
	return &FirestoreService{client: client, log: log}, nil
}

// GetInvoice retrieves an invoice from Firestore
func (s *FirestoreService) GetInvoice(ctx context.Context, invoiceID string) (*models.Invoice, error) {
	s.log.WithField("invoiceID", invoiceID).Info("Getting invoice from Firestore")
	dsnap, err := s.client.Collection("invoices").Doc(invoiceID).Get(ctx)
	if err != nil {
		s.log.WithError(err).WithField("invoiceID", invoiceID).Error("Failed to get invoice from Firestore")
		return nil, err
	}

	var invoice models.Invoice
	if err := dsnap.DataTo(&invoice); err != nil {
		s.log.WithError(err).WithField("invoiceID", invoiceID).Error("Failed to decode invoice data")
		return nil, err
	}

	s.log.WithField("invoiceID", invoiceID).Info("Successfully retrieved invoice")
	return &invoice, nil
}

// GetUser retrieves a user from Firestore
func (s *FirestoreService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	s.log.WithField("userID", userID).Info("Getting user from Firestore")
	dsnap, err := s.client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		s.log.WithError(err).WithField("userID", userID).Error("Failed to get user from Firestore")
		return nil, err
	}

	var user models.User
	if err := dsnap.DataTo(&user); err != nil {
		s.log.WithError(err).WithField("userID", userID).Error("Failed to decode user data")
		return nil, err
	}

	s.log.WithField("userID", userID).Info("Successfully retrieved user")
	return &user, nil
}

// UpdateInvoiceStatus updates the status of an invoice
func (s *FirestoreService) UpdateInvoiceStatus(ctx context.Context, invoiceID, status string) error {
	s.log.WithFields(logrus.Fields{"invoiceID": invoiceID, "status": status}).Info("Updating invoice status in Firestore")
	_, err := s.client.Collection("invoices").Doc(invoiceID).Update(ctx, []firestore.Update{
		{Path: "status", Value: status},
	})
	if err != nil {
		s.log.WithError(err).WithFields(logrus.Fields{"invoiceID": invoiceID, "status": status}).Error("Failed to update invoice status")
		return err
	}

	s.log.WithFields(logrus.Fields{"invoiceID": invoiceID, "status": status}).Info("Successfully updated invoice status")
	return nil
}
