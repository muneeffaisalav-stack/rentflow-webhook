package models

// Button represents a button in a WhatsApp message
type Button struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	ID    string `json:"id"`
}
