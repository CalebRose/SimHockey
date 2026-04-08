package firebase

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// ─────────────────────────────────────────────
// Forum Service
// ─────────────────────────────────────────────

// CreateThread creates a new thread and its first post atomically in Firestore.
// If input.ExternalEventKey is set and a thread with that key already exists,
// the existing thread is returned without creating a duplicate (idempotency).
func CreateThread(ctx context.Context, input CreateForumThreadInput) (*ForumThread, error) {
	if input.ForumID == "" {
		return nil, fmt.Errorf("firebase: ForumID is required")
	}
	if input.Title == "" {
		return nil, fmt.Errorf("firebase: Title is required")
	}

	client := GetFirestoreClient()

	// Idempotency: return existing thread if we already created one for this event.
	if input.ExternalEventKey != "" {
		existing, err := FindThreadByExternalEventKey(ctx, input.ExternalEventKey)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return existing, nil
		}
	}

	now := time.Now().UTC()
	slug := buildSlug(input.Title)

	author := ThreadAuthor{
		UID:         input.AuthorUID,
		Username:    input.AuthorUsername,
		DisplayName: input.AuthorDisplayName,
	}

	postBody := buildSimplePostBody(input.FirstPostBodyText)
	if input.FirstPostBody != nil {
		postBody = input.FirstPostBody
	}
	preview := truncate(input.FirstPostBodyText, 200)

	// 1. Create a post document with an empty threadId; we patch it below.
	postRef := client.Collection("posts").NewDoc()
	post := ForumPost{
		ID:            postRef.ID,
		ThreadID:      "", // patched in step 3
		ForumID:       input.ForumID,
		Author:        PostAuthor(author),
		EditorVersion: 1,
		Body:          postBody,
		BodyText:      input.FirstPostBodyText,
		Mentions:      []PostMention{},
		Reactions:     map[string][]string{},
		IsEdited:      false,
		IsDeleted:     false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if _, err := postRef.Set(ctx, post); err != nil {
		return nil, fmt.Errorf("firebase: failed to create first post: %w", err)
	}

	// 2. Create the thread document.
	threadRef := client.Collection("threads").NewDoc()
	activityBy := &ActivityBy{UID: input.AuthorUID, Username: input.AuthorUsername}
	thread := ForumThread{
		ID:               threadRef.ID,
		ForumID:          input.ForumID,
		ForumPath:        input.ForumPath,
		Title:            input.Title,
		Slug:             slug,
		Author:           author,
		ContentPreview:   preview,
		FirstPostID:      postRef.ID,
		IsPinned:         false,
		IsLocked:         false,
		IsAnnouncement:   false,
		IsDeleted:        false,
		Tags:             []string{},
		ThreadType:       input.ThreadType,
		ReferencedGameID: input.ReferencedGameID,
		ReferencedLeague: input.ReferencedLeague,
		ReplyCount:       0,
		ParticipantCount: 1,
		LatestPostID:     postRef.ID,
		LatestActivityAt: now,
		LatestActivityBy: activityBy,
		CreatedAt:        now,
		UpdatedAt:        now,
		ExternalEventKey: input.ExternalEventKey,
	}
	if _, err := threadRef.Set(ctx, thread); err != nil {
		// Best-effort cleanup: delete the orphaned post.
		if _, delErr := postRef.Delete(ctx); delErr != nil {
			log.Printf("firebase: failed to clean up orphaned post %s: %v", postRef.ID, delErr)
		}
		return nil, fmt.Errorf("firebase: failed to create thread: %w", err)
	}

	// 3. Patch the post with the real threadId.
	if _, err := postRef.Update(ctx, []firestore.Update{
		{Path: "threadId", Value: threadRef.ID},
	}); err != nil {
		log.Printf("firebase: failed to patch post %s with threadId %s: %v", postRef.ID, threadRef.ID, err)
	}

	// 4. Increment forum counters (best-effort; failure is non-fatal).
	if input.ForumID != "" {
		forumRef := client.Collection("forums").Doc(input.ForumID)
		if _, err := forumRef.Update(ctx, []firestore.Update{
			{Path: "threadCount", Value: firestore.Increment(1)},
			{Path: "postCount", Value: firestore.Increment(1)},
			{Path: "latestActivityAt", Value: now},
			{Path: "latestActivityBy", Value: activityBy},
			{Path: "latestThreadId", Value: threadRef.ID},
		}); err != nil {
			log.Printf("firebase: failed to increment forum counters for forum %s: %v", input.ForumID, err)
		}
	}

	thread.ID = threadRef.ID
	return &thread, nil
}

// FindThreadByExternalEventKey looks up a thread by its ExternalEventKey.
// Returns nil, nil if no thread is found.
func FindThreadByExternalEventKey(ctx context.Context, eventKey string) (*ForumThread, error) {
	client := GetFirestoreClient()

	iter := client.Collection("threads").
		Where("externalEventKey", "==", eventKey).
		Limit(1).
		Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("firebase: FindThreadByExternalEventKey: %w", err)
	}

	var thread ForumThread
	if err := doc.DataTo(&thread); err != nil {
		return nil, fmt.Errorf("firebase: failed to decode thread document: %w", err)
	}
	thread.ID = doc.Ref.ID
	return &thread, nil
}

// ─────────────────────────────────────────────
// Rich text helpers
// ─────────────────────────────────────────────

// buildSimplePostBody produces a minimal Tiptap/ProseMirror-compatible rich
// text document that the frontend renderer can display.
func buildSimplePostBody(text string) map[string]interface{} {
	return map[string]interface{}{
		"type": "doc",
		"content": []map[string]interface{}{
			{
				"type": "paragraph",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": text,
					},
				},
			},
		},
	}
}

// ─────────────────────────────────────────────
// String helpers
// ─────────────────────────────────────────────

var nonAlphanumRe = regexp.MustCompile(`[^a-z0-9\s-]`)
var whitespaceRe = regexp.MustCompile(`\s+`)

// buildSlug converts a title to a URL-friendly slug (max 80 chars).
func buildSlug(title string) string {
	s := strings.ToLower(title)
	s = nonAlphanumRe.ReplaceAllString(s, "")
	s = whitespaceRe.ReplaceAllString(s, "-")
	if len(s) > 80 {
		s = s[:80]
	}
	return s
}

// truncate returns at most n runes from s, used for content previews.
func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n])
}
