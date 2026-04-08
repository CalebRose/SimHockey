package firebase

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// ResolveUIDsByUsernames queries the Firestore "users" collection for each
// username and returns the matching document IDs (i.e. Firebase Auth UIDs).
// Usernames that are not found are silently skipped and a warning is logged.
func ResolveUIDsByUsernames(ctx context.Context, usernames []string) []string {
	if len(usernames) == 0 {
		return nil
	}

	client := GetFirestoreClient()
	uids := make([]string, 0, len(usernames))

	for _, username := range usernames {
		if username == "" {
			continue
		}
		uid, err := resolveUID(ctx, client, username)
		if err != nil {
			log.Printf("firebase: could not resolve UID for username %q: %v", username, err)
			continue
		}
		uids = append(uids, uid)
	}

	return uids
}

func resolveUID(ctx context.Context, client *firestore.Client, username string) (string, error) {
	iter := client.Collection("users").
		Where("username", "==", username).
		Limit(1).
		Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return doc.Ref.ID, nil
}
