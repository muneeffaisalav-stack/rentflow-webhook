package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"rentflow-webhook/internal/models"
)

// WhatsappService is a service for interacting with the WhatsApp API
type WhatsappService struct {
	accessToken string
	phoneNumberID string
	logger      *logrus.Logger
}

// NewWhatsappService creates a new WhatsappService
func NewWhatsappService(accessToken, phoneNumberID string, logger *logrus.Logger) *WhatsappService {
	return &WhatsappService{
		accessToken: accessToken,
		phoneNumberID: phoneNumberID,
		logger:      logger,
	}
}

// SendButtonMessage sends a message with buttons to a user
func (s *WhatsappService) SendButtonMessage(to, text string, buttons []models.Button) error {
	url := fmt.Sprintf("https://graph.facebook.com/v13.0/%s/messages", s.phoneNumberID)

	action := models.Action{Buttons: buttons}
	body := &models.ButtonMessage{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "interactive",
		Interactive:      models.Interactive{Type: "button", Body: models.Body{Text: text}, Action: action},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// SendTemplateMessage sends a template message to a user
func (s *WhatsappService) SendTemplateMessage(to, templateName, languageCode string, components []map[string]interface{}) error {
	url := fmt.Sprintf("https://graph.facebook.com/v13.0/%s/messages", s.phoneNumberID)

	// Create the template message body
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "template",
		"template": map[string]interface{}{
			"name":     templateName,
			"language": map[string]string{"code": languageCode},
			"components": components,
		},
	}

	// Marshal the body to JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		s.logger.WithError(err).Error("Failed to marshal template message body")
		return err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		s.logger.WithError(err).Error("Failed to create new HTTP request")
		return err
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.WithError(err).Error("Failed to send template message")
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		var responseBody map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
			s.logger.WithError(err).Warn("Failed to decode error response body")
		}
		s.logger.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"response_body": responseBody,
		}).Error("Received non-OK status from WhatsApp API")
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	s.logger.Info("Template message sent successfully")
	return nil
}
