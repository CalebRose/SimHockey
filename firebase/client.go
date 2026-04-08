package firebase

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var (
	once            sync.Once
	firestoreClient *firestore.Client
)

// GetFirestoreClient returns a singleton Firestore client.
//
// Credential resolution order:
//  1. Individual FIREBASE_* environment variables (recommended for production).
//  2. FIREBASE_SERVICE_ACCOUNT_KEY env var treated as a file path.
//  3. secrets/serviceAccountKey.json relative to the working directory (local dev).
func GetFirestoreClient() *firestore.Client {
	once.Do(func() {
		ctx := context.Background()

		var credOpt option.ClientOption
		if jsonBytes := buildCredentialsFromEnv(); jsonBytes != nil {
			credOpt = option.WithCredentialsJSON(jsonBytes)
		} else {
			credOpt = option.WithCredentialsFile(resolveKeyPath())
		}

		app, err := firebase.NewApp(ctx, nil, credOpt)
		if err != nil {
			log.Fatalf("firebase: failed to initialise app: %v", err)
		}
		client, err := app.Firestore(ctx)
		if err != nil {
			log.Fatalf("firebase: failed to create Firestore client: %v", err)
		}
		firestoreClient = client
	})
	return firestoreClient
}

// buildCredentialsFromEnv assembles a service-account JSON document from
// individual environment variables.  Returns nil if FIREBASE_PRIVATE_KEY is
// not set, signalling that the file-based fallback should be used instead.
//
// Required env vars:
//
//	FIREBASE_PROJECT_ID
//	FIREBASE_PRIVATE_KEY_ID
//	FIREBASE_PRIVATE_KEY      (PEM block; literal \n sequences are expanded automatically)
//	FIREBASE_CLIENT_EMAIL
//	FIREBASE_CLIENT_ID
//	FIREBASE_CLIENT_X509_CERT_URL
func buildCredentialsFromEnv() []byte {
	privateKey := os.Getenv("FIREBASE_PRIVATE_KEY")
	if privateKey == "" {
		return nil
	}

	// Environment variables store the PEM block with literal \n instead of
	// real newlines.  The key parser requires actual newline characters.
	privateKey = strings.ReplaceAll(privateKey, `\n`, "\n")

	creds := map[string]string{
		"type":                        "service_account",
		"project_id":                  os.Getenv("FIREBASE_PROJECT_ID"),
		"private_key_id":              os.Getenv("FIREBASE_PRIVATE_KEY_ID"),
		"private_key":                 privateKey,
		"client_email":                os.Getenv("FIREBASE_CLIENT_EMAIL"),
		"client_id":                   os.Getenv("FIREBASE_CLIENT_ID"),
		"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
		"token_uri":                   "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url":        os.Getenv("FIREBASE_CLIENT_X509_CERT_URL"),
	}

	b, err := json.Marshal(creds)
	if err != nil {
		log.Fatalf("firebase: failed to marshal credentials from env: %v", err)
	}
	return b
}

// resolveKeyPath returns the path to the service account key file.
// Override with FIREBASE_SERVICE_ACCOUNT_KEY for a non-default path.
func resolveKeyPath() string {
	if path := os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY"); path != "" {
		return path
	}
	return "secrets/serviceAccountKey.json"
}
