package models

import "cloud.google.com/go/firestore"

// Invoice represents an invoice in Firestore
type Invoice struct {
	Amount       float64                `firestore:"amount"`
	DueDate      string                 `firestore:"due_date"`
	Month        string                 `firestore:"month"`
	Status       string                 `firestore:"status"`
	PropertyID   string                 `firestore:"propertyId"` // Changed to match DB
	TenantID     string                 `firestore:"tenantId"`   // Changed to match DB
	LandlordID   string                 `firestore:"landlordId"` // Changed to match DB
	PropertyName string                 `firestore:"property_name"`
	// Deprecated reference fields
	TenantRef   *firestore.DocumentRef `firestore:"tenant_ref,omitempty"`
	LandlordRef *firestore.DocumentRef `firestore:"landlord_ref,omitempty"`
}
