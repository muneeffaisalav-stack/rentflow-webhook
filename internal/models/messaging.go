package models

// ButtonMessage represents a message with buttons
type ButtonMessage struct {
	MessagingProduct string      `json:"messaging_product"`
	To               string      `json:"to"`
	Type             string      `json:"type"`
	Interactive      Interactive `json:"interactive"`
}

// Interactive represents an interactive component of a message
type Interactive struct {
	Type   string `json:"type"`
	Body   Body   `json:"body"`
	Action Action `json:"action"`
}

// Body represents the body of an interactive message
type Body struct {
	Text string `json:"text"`
}

// Action represents an action in an interactive message
type Action struct {
	Buttons []Button `json:"buttons"`
}
