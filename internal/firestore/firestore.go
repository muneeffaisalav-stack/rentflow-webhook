package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

// Client is a client for interacting with Firestore

type Client struct {
	client *firestore.Client
}

// NewClient creates a new Firestore client

func NewClient(ctx context.Context, projectID, clientEmail, privateKey string) (*Client, error) {
	creds := option.WithCredentialsJSON([]byte(fmt.Sprintf(`{
	  "type": "service_account",
	  "project_id": "%s",
	  "private_key_id": "",
	  "private_key": "%s",
	  "client_email": "%s",
	  "client_id": "",
	  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
	  "token_uri": "https://oauth2.googleapis.com/token",
	  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
	  "client_x509_cert_url": ""
	}`, projectID, privateKey, clientEmail)))

	client, err := firestore.NewClient(ctx, projectID, creds)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client: %w", err)
	}

	return &Client{client: client}, nil
}
