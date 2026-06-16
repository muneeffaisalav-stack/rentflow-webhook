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
	From        string `json:"from"`
	ID          string `json:"id"`
	Timestamp   string `json:"timestamp"`
	Interactive struct {
		Type       string `json:"type"`
		ButtonReply struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"button_reply"`
	} `json:"interactive"`
}
