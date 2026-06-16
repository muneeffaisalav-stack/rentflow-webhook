package models

import "cloud.google.com/go/firestore"

// Invoice represents an invoice in Firestore
type Invoice struct {
	Amount      float64                `firestore:"amount"`
	DueDate     string                 `firestore:"due_date"`
	Status      string                 `firestore:"status"`
	TenantRef   *firestore.DocumentRef `firestore:"tenant_ref"`
	LandlordRef *firestore.DocumentRef `firestore:"landlord_ref"`
}
