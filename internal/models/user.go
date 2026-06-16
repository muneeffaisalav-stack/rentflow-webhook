package models

// User represents a user in Firestore
type User struct {
	Name        string `firestore:"name"`
	PhoneNumber string `firestore:"phone_number"`
}
