package firebase

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// ─────────────────────────────────────────────
// Notification Service
// ─────────────────────────────────────────────

// NotifyTeamInjury notifies a team's coaches or owners that a player was injured
// during a game.  Idempotent via SourceEventKey (keyed per player per game).
func NotifyTeamInjury(ctx context.Context, input TeamInjuryNotificationInput) error {
	if len(input.RecipientUIDs) == 0 {
		return nil
	}

	daysStr := "1 day"
	if input.DaysOfRecovery > 1 {
		daysStr = fmt.Sprintf("%d days", input.DaysOfRecovery)
	}
	message := fmt.Sprintf(
		"%s (%s) suffered %s and is expected to miss approximately %s. Check the roster for details.",
		input.PlayerName, input.Position, input.InjuryType, daysStr,
	)
	linkTo := BuildTeamRosterRoute(input.League, input.TeamID)

	return writeNotificationsIfNew(ctx, input.RecipientUIDs, ForumNotification{
		Type:           NotificationTypeInjury,
		Domain:         input.Domain,
		LinkTo:         linkTo,
		Message:        message,
		ActorUsername:  "SimSN",
		IsRead:         false,
		SourceEventKey: input.SourceEventKey,
	})
}

// NotifyRecruitSigned creates one notification per recipient when a recruit
// commits to a team.  Idempotent via SourceEventKey.
func NotifyRecruitSigned(ctx context.Context, input RecruitSignedNotificationInput) error {
	if len(input.RecipientUIDs) == 0 {
		return nil
	}

	message := fmt.Sprintf("%s has signed with %s.", input.RecruitName, input.TeamName)
	linkTo := BuildTeamRecruitingRoute(input.League, input.TeamID)

	return writeNotificationsIfNew(ctx, input.RecipientUIDs, ForumNotification{
		Type:           NotificationTypeRecruiting,
		Domain:         input.Domain,
		LinkTo:         linkTo,
		Message:        message,
		ActorUsername:  "SimSN",
		IsRead:         false,
		SourceEventKey: input.SourceEventKey,
	})
}

// NotifyAffiliatePlayerOffer notifies a PHL team's owner/GM that another team
// has placed an offer on one of their affiliate players.  Idempotent via SourceEventKey.
func NotifyAffiliatePlayerOffer(ctx context.Context, input AffiliatePlayerOfferNotificationInput) error {
	if len(input.RecipientUIDs) == 0 {
		return nil
	}

	message := fmt.Sprintf(
		"%s have placed an offer on %s %s to pick up from your affiliate roster.",
		input.OfferingTeam, input.Position, input.PlayerName,
	)
	linkTo := BuildTeamRosterRoute("phl", input.OwnerTeamID)

	return writeNotificationsIfNew(ctx, input.RecipientUIDs, ForumNotification{
		Type:           NotificationTypeFreeAgency,
		Domain:         DomainPHL,
		LinkTo:         linkTo,
		Message:        message,
		ActorUsername:  "SimSN",
		IsRead:         false,
		SourceEventKey: input.SourceEventKey,
	})
}

// NotifyTransferIntention notifies a coach that one of their players has declared
// an intention to enter the transfer portal.  Idempotent via SourceEventKey.
func NotifyTransferIntention(ctx context.Context, input TransferIntentionNotificationInput) error {
	if len(input.RecipientUIDs) == 0 {
		return nil
	}

	message := fmt.Sprintf(
		"%d star %s %s has a %s likeliness of entering the transfer portal. Please navigate to the Roster page to submit a promise.",
		input.Stars, input.Position, input.PlayerName, input.TransferLikeliness,
	)
	linkTo := BuildTeamRosterRoute("chl", input.TeamID)

	return writeNotificationsIfNew(ctx, input.RecipientUIDs, ForumNotification{
		Type:           NotificationTypeTransfer,
		Domain:         DomainCHL,
		LinkTo:         linkTo,
		Message:        message,
		ActorUsername:  "SimSN",
		IsRead:         false,
		SourceEventKey: input.SourceEventKey,
	})
}

// NotifyTransferPortalSigning notifies a coach that a transfer portal player has
// officially signed with their team.  Idempotent via SourceEventKey.
func NotifyTransferPortalSigning(ctx context.Context, input TransferPortalSigningNotificationInput) error {
	if len(input.RecipientUIDs) == 0 {
		return nil
	}

	message := fmt.Sprintf(
		"%d star %s %s (formerly of %s) has signed with %s via the transfer portal.",
		input.Stars, input.Position, input.PlayerName, input.PreviousTeam, input.TeamName,
	)
	linkTo := BuildTeamRosterRoute("chl", input.TeamID)

	return writeNotificationsIfNew(ctx, input.RecipientUIDs, ForumNotification{
		Type:           NotificationTypeTransfer,
		Domain:         DomainCHL,
		LinkTo:         linkTo,
		Message:        message,
		ActorUsername:  "SimSN",
		IsRead:         false,
		SourceEventKey: input.SourceEventKey,
	})
}

// NotifyGameplanIssue sends a depth-chart / gameplan penalty notification to
// the coach or owner of a team.  Idempotent via SourceEventKey.
func NotifyGameplanIssue(ctx context.Context, input GameplanNotificationInput) error {
	if len(input.RecipientUIDs) == 0 {
		return nil
	}
	linkTo := BuildTeamGameplanRoute(input.League, input.TeamID)

	return writeNotificationsIfNew(ctx, input.RecipientUIDs, ForumNotification{
		Type:           NotificationTypeGameplan,
		Domain:         input.Domain,
		LinkTo:         linkTo,
		Message:        input.Message,
		ActorUsername:  "SimSN",
		IsRead:         false,
		SourceEventKey: input.SourceEventKey,
	})
}

// NotifyScheduleEvent notifies a coach about a game-request lifecycle event
// such as acceptance, rejection, or an admin veto. Idempotent via SourceEventKey.
func NotifyScheduleEvent(ctx context.Context, input ScheduleEventNotificationInput) error {
	if len(input.RecipientUIDs) == 0 {
		return nil
	}
	return writeNotificationsIfNew(ctx, input.RecipientUIDs, ForumNotification{
		Type:           NotificationTypeSystem,
		Domain:         input.Domain,
		LinkTo:         BuildTeamRosterRoute(input.League, input.TeamID),
		Message:        input.Message,
		ActorUsername:  "SimSN",
		IsRead:         false,
		SourceEventKey: input.SourceEventKey,
	})
}

// ─────────────────────────────────────────────
// Internal helpers
// ─────────────────────────────────────────────

// writeNotificationsIfNew writes one notification document per recipient UID,
// skipping any recipient that already has a document with the same
// SourceEventKey (idempotency guard).
func writeNotificationsIfNew(
	ctx context.Context,
	recipientUIDs []string,
	template ForumNotification,
) error {
	client := GetFirestoreClient()
	col := client.Collection("notifications")
	now := time.Now().UTC()

	for _, uid := range recipientUIDs {
		if uid == "" {
			continue
		}

		// Idempotency: skip if this event was already delivered to this recipient.
		if template.SourceEventKey != "" {
			exists, err := notificationExists(ctx, col, uid, template.SourceEventKey)
			if err != nil {
				log.Printf("firebase: idempotency check failed for uid=%s key=%s: %v", uid, template.SourceEventKey, err)
			}
			if exists {
				continue
			}
		}

		ref := col.NewDoc()
		n := template
		n.ID = ref.ID
		n.UID = uid
		n.CreatedAt = now

		if _, err := ref.Set(ctx, n); err != nil {
			log.Printf("firebase: failed to write notification for uid=%s: %v", uid, err)
		}
	}

	return nil
}

// notificationExists returns true when a notification doc already exists for
// the given uid and sourceEventKey.
func notificationExists(
	ctx context.Context,
	col *firestore.CollectionRef,
	uid string,
	sourceEventKey string,
) (bool, error) {
	iter := col.
		Where("uid", "==", uid).
		Where("sourceEventKey", "==", sourceEventKey).
		Limit(1).
		Documents(ctx)
	defer iter.Stop()

	_, err := iter.Next()
	if err == iterator.Done {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
