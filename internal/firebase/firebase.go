package firebase

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"rentflow-webhook/internal/config"
)

// FirestoreService is a service for interacting with Firestore
type FirestoreService struct {
	Client *firestore.Client
	log    *logrus.Logger
}

// Invoice represents an invoice in Firestore
type Invoice struct {
	Amount       float64                `firestore:"amount"`
	DueDate      string                 `firestore:"due_date"`
	Month        string                 `firestore:"month"`
	Status       string                 `firestore:"status"`
	PropertyName string                 `firestore:"property_name"`
	LandlordRef  *firestore.DocumentRef `firestore:"landlordRef"`
	TenantRef    *firestore.DocumentRef `firestore:"tenantRef"`
}

// User represents a user in Firestore
type User struct {
	Name        string `firestore:"name"`
	PhoneNumber string `firestore:"phoneNumber"`
}

// NewFirestoreService creates a new FirestoreService
func NewFirestoreService(cfg *config.Config, log *logrus.Logger) (*FirestoreService, error) {
	ctx := context.Background()

	// Set up the Firebase credentials
	opt := option.WithCredentialsJSON([]byte(cfg.FirebaseCredentialsJSON))

	// Initialize the Firebase app
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.WithError(err).Error("Failed to create Firebase app")
		return nil, err
	}

	// Get a Firestore client
	client, err := app.Firestore(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to create Firestore client")
		return nil, err
	}

	log.Info("Firestore service created successfully")

	return &FirestoreService{
		Client: client,
		log:    log,
	}, nil
}

// GetInvoice retrieves an invoice from Firestore
func (s *FirestoreService) GetInvoice(ctx context.Context, invoiceID string) (*Invoice, error) {
	s.log.WithField("invoiceID", invoiceID).Info("Getting invoice from Firestore")
	doc, err := s.Client.Collection("invoices").Doc(invoiceID).Get(ctx)
	if err != nil {
		return nil, err
	}
	var invoice Invoice
	if err := doc.DataTo(&invoice); err != nil {
		return nil, err
	}
	return &invoice, nil
}

// GetUser retrieves a user from Firestore
func (s *FirestoreService) GetUser(ctx context.Context, userID string) (*User, error) {
	s.log.WithField("userID", userID).Info("Getting user from Firestore")
	doc, err := s.Client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		return nil, err
	}
	var user User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateInvoiceStatus updates the status of an invoice in Firestore
func (s *FirestoreService) UpdateInvoiceStatus(ctx context.Context, invoiceID, status string) error {
	s.log.WithFields(logrus.Fields{
		"invoiceID": invoiceID,
		"status":    status,
	}).Info("Updating invoice status in Firestore")
	_, err := s.Client.Collection("invoices").Doc(invoiceID).Update(ctx, []firestore.Update{
		{Path: "status", Value: status},
	})
	return err
}
