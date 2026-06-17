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

// TemplateMessage represents a template message
type TemplateMessage struct {
	MessagingProduct string   `json:"messaging_product"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Template         Template `json:"template"`
}

// Template represents a template in a message
type Template struct {
	Name       string      `json:"name"`
	Language   Language    `json:"language"`
	Components []Component `json:"components"`
}

// Language represents the language of a template
type Language struct {
	Code string `json:"code"`
}

// Component represents a component of a template
type Component struct {
	Type       string      `json:"type"`
	SubType    string      `json:"sub_type,omitempty"`
	Index      string      `json:"index,omitempty"`
	Parameters []Parameter `json:"parameters"`
}

// Parameter represents a parameter in a component
type Parameter struct {
	Type    string `json:"type"`
	Text    string `json:"text,omitempty"`
	Payload string `json:"payload,omitempty"`
}
