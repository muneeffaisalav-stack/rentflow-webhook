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
	From        string       `json:"from"`
	ID          string       `json:"id"`
	Timestamp   string       `json:"timestamp"`
	Type        string       `json:"type"`
	Button      *Button      `json:"button,omitempty"`
	Interactive *Interactive `json:"interactive,omitempty"`
}

// Button represents a button in a WhatsApp message
type Button struct {
	Payload string `json:"payload"`
	Text    string `json:"text"`
}

// Interactive represents an interactive message in WhatsApp
type Interactive struct {
	Type        string      `json:"type"`
	ButtonReply ButtonReply `json:"button_reply"`
}

// ButtonReply is the reply from a button in an interactive message
type ButtonReply struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}
