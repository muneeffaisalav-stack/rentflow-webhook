package models

// WhatsappWebhook is the top-level struct for a WhatsApp webhook
type WhatsappWebhook struct {
	Entry []struct {
		Changes []struct {
			Value struct {
				Messages []Message `json:"messages"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}

// Message is a single message from a WhatsApp webhook
type Message struct {
	From        string               `json:"from"`
	ID          string               `json:"id"`
	Timestamp   string               `json:"timestamp"`
	Type        string               `json:"type"`
	Button      *Button              `json:"button,omitempty"`
	Interactive *IncomingInteractive `json:"interactive,omitempty"`
}

// IncomingInteractive represents the interactive part of an incoming message from WhatsApp
type IncomingInteractive struct {
	Type        string      `json:"type"`
	ButtonReply ButtonReply `json:"button_reply"`
}

// ButtonReply represents a button reply from WhatsApp
type ButtonReply struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}
