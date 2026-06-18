package repository

import (
	"context"
	"rentflow-webhook/internal/models"
)

// Repository defines the interface for interacting with the database

type Repository interface {
	GetInvoiceByID(ctx context.Context, id string) (*models.Invoice, error)
	UpdateInvoice(ctx context.Context, invoice *models.Invoice) error
	UpdateInvoiceStatus(ctx context.Context, id, status string) error
	SaveUTR(ctx context.Context, id, utr string) error
	StoreVerificationToken(ctx context.Context, id, token string) error
	ClearVerificationToken(ctx context.Context, id string) error
	CreateClaimSession(ctx context.Context, session *models.PaymentClaimSession) error
	GetClaimSession(ctx context.Context, tenantPhone string) (*models.PaymentClaimSession, error)
	DeleteClaimSession(ctx context.Context, tenantPhone string) error
	GetTenant(ctx context.Context, id string) (*models.Tenant, error)
	GetOwner(ctx context.Context, id string) (*models.Owner, error)
}
