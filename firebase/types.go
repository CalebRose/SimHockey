package firebase

import "time"

// ─────────────────────────────────────────────
// Shared sub-types
// ─────────────────────────────────────────────

// ThreadAuthor mirrors the frontend ThreadAuthor shape stored in Firestore.
type ThreadAuthor struct {
	UID         string `firestore:"uid"`
	Username    string `firestore:"username"`
	DisplayName string `firestore:"displayName"`
	LogoURL     string `firestore:"logoUrl,omitempty"`
}

// PostAuthor mirrors the frontend PostAuthor shape stored in Firestore.
type PostAuthor struct {
	UID         string `firestore:"uid"`
	Username    string `firestore:"username"`
	DisplayName string `firestore:"displayName"`
	LogoURL     string `firestore:"logoUrl,omitempty"`
}

// ActivityBy is the compact author snapshot used for "latest activity" fields.
type ActivityBy struct {
	UID      string `firestore:"uid"`
	Username string `firestore:"username"`
}

// PostMention is a user mentioned inside a post body.
type PostMention struct {
	UID      string `firestore:"uid"`
	Username string `firestore:"username"`
}

// ─────────────────────────────────────────────
// Forum Thread
// ─────────────────────────────────────────────

// Thread type constants (must stay in sync with the frontend ThreadType union).
const (
	ThreadTypeStandard      = "standard"
	ThreadTypeGameReference = "game_reference"
	ThreadTypePoll          = "poll"
)

// CreatedByType constants.
const (
	CreatedByUser   = "user"
	CreatedBySystem = "system"
	CreatedByBot    = "bot"
)

// ForumThread is the Firestore document shape for the "threads" collection.
type ForumThread struct {
	ID               string       `firestore:"id"`
	ForumID          string       `firestore:"forumId"`
	ForumPath        []string     `firestore:"forumPath"`
	Title            string       `firestore:"title"`
	Slug             string       `firestore:"slug"`
	Author           ThreadAuthor `firestore:"author"`
	ContentPreview   string       `firestore:"contentPreview"`
	FeatureImageURL  string       `firestore:"featureImageUrl,omitempty"`
	FirstPostID      string       `firestore:"firstPostId"`
	IsPinned         bool         `firestore:"isPinned"`
	IsLocked         bool         `firestore:"isLocked"`
	IsAnnouncement   bool         `firestore:"isAnnouncement"`
	IsDeleted        bool         `firestore:"isDeleted"`
	Tags             []string     `firestore:"tags"`
	ThreadType       string       `firestore:"threadType"`
	PollID           string       `firestore:"pollId,omitempty"`
	ReferencedGameID string       `firestore:"referencedGameId,omitempty"`
	ReferencedLeague string       `firestore:"referencedLeague,omitempty"`
	ReplyCount       int          `firestore:"replyCount"`
	ParticipantCount int          `firestore:"participantCount"`
	LatestPostID     string       `firestore:"latestPostId,omitempty"`
	LatestActivityAt time.Time    `firestore:"latestActivityAt"`
	LatestActivityBy *ActivityBy  `firestore:"latestActivityBy,omitempty"`
	CreatedAt        time.Time    `firestore:"createdAt"`
	UpdatedAt        time.Time    `firestore:"updatedAt"`

	// ExternalEventKey is set on system-generated threads for idempotency.
	ExternalEventKey string `firestore:"externalEventKey,omitempty"`
}

// ─────────────────────────────────────────────
// Forum Post
// ─────────────────────────────────────────────

// ForumPost is the Firestore document shape for the "posts" collection.
type ForumPost struct {
	ID               string                 `firestore:"id"`
	ThreadID         string                 `firestore:"threadId"`
	ForumID          string                 `firestore:"forumId"`
	Author           PostAuthor             `firestore:"author"`
	EditorVersion    int                    `firestore:"editorVersion"`
	Body             map[string]interface{} `firestore:"body"`
	BodyText         string                 `firestore:"bodyText"`
	QuotedPostID     string                 `firestore:"quotedPostId,omitempty"`
	ReplyToPostID    string                 `firestore:"replyToPostId,omitempty"`
	Mentions         []PostMention          `firestore:"mentions"`
	Reactions        map[string][]string    `firestore:"reactions"`
	IsEdited         bool                   `firestore:"isEdited"`
	EditedAt         *time.Time             `firestore:"editedAt,omitempty"`
	EditedBy         string                 `firestore:"editedBy,omitempty"`
	IsDeleted        bool                   `firestore:"isDeleted"`
	DeletedAt        *time.Time             `firestore:"deletedAt,omitempty"`
	DeletedBy        string                 `firestore:"deletedBy,omitempty"`
	ModerationReason string                 `firestore:"moderationReason,omitempty"`
	CreatedAt        time.Time              `firestore:"createdAt"`
	UpdatedAt        time.Time              `firestore:"updatedAt"`
}

// ─────────────────────────────────────────────
// Notification
// ─────────────────────────────────────────────

// Notification type constants.
const (
	NotificationTypeInjury       = "injury"
	NotificationTypeRecruiting   = "recruiting"
	NotificationTypeGameplan     = "gameplan"
	NotificationTypeTrade        = "trade"
	NotificationTypeDraft        = "draft"
	NotificationTypeFreeAgency   = "free_agency"
	NotificationTypeTransfer     = "transfer"
	NotificationTypeSystem       = "system"
	NotificationTypeForumReply   = "reply"
	NotificationTypeForumMention = "mention"
)

// Notification domain constants.
const (
	DomainCHL    = "chl"
	DomainPHL    = "phl"
	DomainForum  = "forum"
	DomainSystem = "system"
)

// ForumNotification is the Firestore document shape for the "notifications" collection.
type ForumNotification struct {
	ID            string    `firestore:"id"`
	UID           string    `firestore:"uid"`    // Firebase Auth UID of the recipient
	Type          string    `firestore:"type"`   // NotificationForumType
	Domain        string    `firestore:"domain"` // NotificationDomain
	LinkTo        string    `firestore:"linkTo,omitempty"`
	ThreadID      string    `firestore:"threadId,omitempty"`
	PostID        string    `firestore:"postId,omitempty"`
	ActorUID      string    `firestore:"actorUid,omitempty"`
	ActorUsername string    `firestore:"actorUsername,omitempty"`
	Message       string    `firestore:"message"`
	IsRead        bool      `firestore:"isRead"`
	CreatedAt     time.Time `firestore:"createdAt"`

	// SourceEventKey supports idempotency checks.
	SourceEventKey string `firestore:"sourceEventKey,omitempty"`
}

// ─────────────────────────────────────────────
// Service input types
// ─────────────────────────────────────────────

// CreateForumThreadInput carries all the data required to create a thread + its first post atomically.
type CreateForumThreadInput struct {
	ForumID           string
	ForumPath         []string
	Title             string
	AuthorUID         string
	AuthorUsername    string
	AuthorDisplayName string
	CreatedByType     string
	ThreadType        string
	FirstPostBodyText string
	FirstPostBody     map[string]interface{}
	ReferencedGameID  string
	ReferencedLeague  string
	ExternalEventKey  string
	Metadata          map[string]interface{}
}

// TeamInjuryNotificationInput carries the context needed to notify a team's
// coaches or owners that a player was injured during a game.
type TeamInjuryNotificationInput struct {
	League         string
	Domain         string // e.g. DomainCHL, DomainPHL
	TeamID         uint
	TeamName       string
	PlayerID       uint
	PlayerName     string
	Position       string
	InjuryType     string
	DaysOfRecovery int
	GameID         string
	RecipientUIDs  []string
	SourceEventKey string
}

// RecruitSignedNotificationInput carries the context needed to build recruit-signing notifications.
type RecruitSignedNotificationInput struct {
	League         string
	Domain         string // e.g. DomainCHL
	TeamID         uint
	TeamName       string
	RecruitID      uint
	RecruitName    string
	RecipientUIDs  []string
	SourceEventKey string
}

// AffiliatePlayerOfferNotificationInput carries the context needed to notify a PHL
// team's staff that another team has placed an offer on their affiliate player.
type AffiliatePlayerOfferNotificationInput struct {
	OwnerTeamID    uint
	OwnerTeamName  string
	OfferingTeam   string
	PlayerID       uint
	PlayerName     string
	Position       string
	RecipientUIDs  []string
	SourceEventKey string
}

// TransferIntentionNotificationInput carries the context needed to notify a coach
// that one of their players has declared an intention to enter the transfer portal.
type TransferIntentionNotificationInput struct {
	TeamID             uint
	TeamAbbr           string
	PlayerID           uint
	PlayerName         string
	Position           string
	Stars              int
	TransferLikeliness string
	RecipientUIDs      []string
	SourceEventKey     string
}

// TransferPortalSigningNotificationInput carries the context needed to notify a coach
// that a transfer portal player has signed with their team.
type TransferPortalSigningNotificationInput struct {
	TeamID         uint
	TeamName       string
	PlayerID       uint
	PlayerName     string
	Position       string
	Stars          int
	PreviousTeam   string
	RecipientUIDs  []string
	SourceEventKey string
}

// GameplanNotificationInput carries the context needed to notify a coach or owner
// about a depth-chart or gameplan issue.
type GameplanNotificationInput struct {
	League         string
	Domain         string
	TeamID         uint
	TeamName       string
	TeamAbbr       string
	Message        string
	RecipientUIDs  []string
	SourceEventKey string
}

// ScheduleEventNotificationInput carries the context needed to notify a coach
// about a game-request lifecycle event (accepted, rejected, vetoed).
type ScheduleEventNotificationInput struct {
	League         string
	Domain         string // e.g. DomainCHL, DomainPHL
	TeamID         uint
	Message        string // fully-formed message from the caller
	RecipientUIDs  []string
	SourceEventKey string
}
