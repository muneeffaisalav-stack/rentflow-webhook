package repository

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"rentflow-webhook/internal/models"
)

// firestoreRepository is a repository that uses Firestore as a backend

type firestoreRepository struct {
	client *firestore.Client
}

// NewFirestoreRepository creates a new Firestore repository

func NewFirestoreRepository(client *firestore.Client) Repository {
	return &firestoreRepository{client: client}
}

// GetInvoiceByID returns an invoice by its ID

func (r *firestoreRepository) GetInvoiceByID(ctx context.Context, id string) (*models.Invoice, error) {
	log.Printf("Getting invoice by ID: %s", id)
	dsnap, err := r.client.Collection("invoices").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			log.Printf("Invoice with ID %s not found", id)
			return nil, nil
		}
		log.Printf("Error getting invoice with ID %s: %v", id, err)
		return nil, err
	}

	var invoice models.Invoice
	if err := dsnap.DataTo(&invoice); err != nil {
		log.Printf("Error converting invoice data to struct: %v", err)
		return nil, err
	}

	log.Printf("Successfully got invoice with ID: %s", id)
	return &invoice, nil
}

// UpdateInvoice updates an invoice

func (r *firestoreRepository) UpdateInvoice(ctx context.Context, invoice *models.Invoice) error {
	log.Printf("Updating invoice with ID: %s", invoice.ID)
	_, err := r.client.Collection("invoices").Doc(invoice.ID).Set(ctx, invoice)
	if err != nil {
		log.Printf("Error updating invoice with ID %s: %v", invoice.ID, err)
		return err
	}
	log.Printf("Successfully updated invoice with ID: %s", invoice.ID)
	return nil
}

// UpdateInvoiceStatus updates the status of an invoice

func (r *firestoreRepository) UpdateInvoiceStatus(ctx context.Context, id, status string) error {
	log.Printf("Updating status of invoice with ID %s to %s", id, status)
	_, err := r.client.Collection("invoices").Doc(id).Update(ctx, []firestore.Update{
		{Path: "status", Value: status},
	})
	if err != nil {
		log.Printf("Error updating status of invoice with ID %s: %v", id, err)
		return err
	}
	log.Printf("Successfully updated status of invoice with ID %s", id)
	return err
}

// SaveUTR saves the UTR of an invoice

func (r *firestoreRepository) SaveUTR(ctx context.Context, id, utr string) error {
	log.Printf("Saving UTR for invoice with ID %s", id)
	_, err := r.client.Collection("invoices").Doc(id).Update(ctx, []firestore.Update{
		{Path: "utr", Value: utr},
	})
	if err != nil {
		log.Printf("Error saving UTR for invoice with ID %s: %v", id, err)
		return err
	}
	log.Printf("Successfully saved UTR for invoice with ID %s", id)
	return err
}

// StoreVerificationToken stores a verification token for an invoice

func (r *firestoreRepository) StoreVerificationToken(ctx context.Context, id, token string) error {
	log.Printf("Storing verification token for invoice with ID %s", id)
	expiresAt := time.Now().Add(24 * time.Hour)
	_, err := r.client.Collection("invoices").Doc(id).Update(ctx, []firestore.Update{
		{Path: "verificationToken", Value: token},
		{Path: "tokenExpiresAt", Value: expiresAt},
	})
	if err != nil {
		log.Printf("Error storing verification token for invoice with ID %s: %v", id, err)
		return err
	}
	log.Printf("Successfully stored verification token for invoice with ID %s", id)
	return err
}

// ClearVerificationToken clears a verification token for an invoice

func (r *firestoreRepository) ClearVerificationToken(ctx context.Context, id string) error {
	log.Printf("Clearing verification token for invoice with ID %s", id)
	_, err := r.client.Collection("invoices").Doc(id).Update(ctx, []firestore.Update{
		{Path: "verificationToken", Value: firestore.Delete},
		{Path: "tokenExpiresAt", Value: firestore.Delete},
	})
	if err != nil {
		log.Printf("Error clearing verification token for invoice with ID %s: %v", id, err)
		return err
	}
	log.Printf("Successfully cleared verification token for invoice with ID %s", id)
	return err
}

// CreateClaimSession creates a payment claim session

func (r *firestoreRepository) CreateClaimSession(ctx context.Context, session *models.PaymentClaimSession) error {
	log.Printf("Creating claim session for tenant phone: %s", session.TenantPhone)
	_, err := r.client.Collection("payment_claim_sessions").Doc(session.TenantPhone).Set(ctx, session)
	if err != nil {
		log.Printf("Error creating claim session for tenant phone %s: %v", session.TenantPhone, err)
		return err
	}
	log.Printf("Successfully created claim session for tenant phone: %s", session.TenantPhone)
	return nil
}

// GetClaimSession returns a payment claim session by tenant phone

func (r *firestoreRepository) GetClaimSession(ctx context.Context, tenantPhone string) (*models.PaymentClaimSession, error) {
	log.Printf("Getting claim session for tenant phone: %s", tenantPhone)
	dsnap, err := r.client.Collection("payment_claim_sessions").Doc(tenantPhone).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			log.Printf("Claim session for tenant phone %s not found", tenantPhone)
			return nil, nil
		}
		log.Printf("Error getting claim session for tenant phone %s: %v", tenantPhone, err)
		return nil, err
	}

	var session models.PaymentClaimSession
	if err := dsnap.DataTo(&session); err != nil {
		log.Printf("Error converting claim session data to struct: %v", err)
		return nil, err
	}

	log.Printf("Successfully got claim session for tenant phone: %s", tenantPhone)
	return &session, nil
}

// DeleteClaimSession deletes a payment claim session

func (r *firestoreRepository) DeleteClaimSession(ctx context.-Context, tenantPhone string) error {
	log.Printf("Deleting claim session for tenant phone: %s", tenantPhone)
	_, err := r.client.Collection("payment_claim_sessions").Doc(tenantPhone).Delete(ctx)
	if err != nil {
		log.Printf("Error deleting claim session for tenant phone %s: %v", tenantPhone, err)
		return err
	}
	log.Printf("Successfully deleted claim session for tenant phone: %s", tenantPhone)
	return nil
}

// GetTenant returns a tenant by its ID

func (r *firestoreRepository) GetTenant(ctx context.Context, id string) (*models.Tenant, error) {
	log.Printf("Getting tenant by ID: %s", id)
	dsnap, err := r.client.Collection("tenants").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			log.Printf("Tenant with ID %s not found", id)
			return nil, nil
		}
		log.Printf("Error getting tenant with ID %s: %v", id, err)
		return nil, err
	}

	var tenant models.Tenant
	if err := dsnap.DataTo(&tenant); err != nil {
		log.Printf("Error converting tenant data to struct: %v", err)
		return nil, err
	}

	log.Printf("Successfully got tenant with ID: %s", id)
	return &tenant, nil
}

// GetOwner returns an owner by its ID

func (r *firestoreRepository) GetOwner(ctx context.Context, id string) (*models.Owner, error) {
	log.Printf("Getting owner by ID: %s", id)
	dsnap, err := r.client.Collection("owners").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			log.Printf("Owner with ID %s not found", id)
			return nil, nil
		}
		log.Printf("Error getting owner with ID %s: %v", id, err)
		return nil, err
	}

	var owner models.Owner
	if err := dsnap.DataTo(&owner); err != nil {
		log.Printf("Error converting owner data to struct: %v", err)
		return nil, err
	}

	log.Printf("Successfully got owner with ID: %s", id)
	return &owner, nil
}
