package models

// WhatsappWebhook represents the entire incoming webhook payload from WhatsApp.
type WhatsappWebhook struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

// Entry represents a single entry in the webhook payload.
type Entry struct {
	ID      string   `json:"id"`
	Changes []Change `json:"changes"`
}

// Change represents a change object within an entry.
type Change struct {
	Field string `json:"field"`
	Value Value  `json:"value"`
}

// Value contains the details of the change, including metadata, contacts, and messages.
type Value struct {
	MessagingProduct string    `json:"messaging_product"`
	Metadata         Metadata  `json:"metadata"`
	Contacts         []Contact `json:"contacts"`
	Messages         []Message `json:"messages"`
}

// Metadata contains information about the phone number receiving the message.
type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

// Contact contains information about the user who sent the message.
type Contact struct {
	Profile Profile `json:"profile"`
	WaID    string  `json:"wa_id"`
}

// Profile contains the name of the user.
type Profile struct {
	Name string `json:"name"`
}

// Message is a single message from a WhatsApp webhook
type Message struct {
	From        string               `json:"from"`
	ID          string               `json:"id"`
	Timestamp   string               `json:"timestamp"`
	Type        string               `json:"type"`
	Button      *Button              `json:"button,omitempty"`
	Interactive *IncomingInteractive `json:"interactive,omitempty"`
	Context     *Context             `json:"context,omitempty"`
}

// Context contains the context of a replied-to message.
type Context struct {
	From string `json:"from"`
	ID   string `json:"id"`
}

// Button represents a button click from a WhatsApp message.
type Button struct {
	Payload string `json:"payload"`
	Text    string `json:"text"`
}

// IncomingInteractive represents an interactive message reply.
type IncomingInteractive struct {
	Type        string      `json:"type"`
	ButtonReply ButtonReply `json:"button_reply"`
}

// ButtonReply represents the user's reply to a button message.
type ButtonReply struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}
