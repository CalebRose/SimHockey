package firebase

import (
	"fmt"
	"strconv"
)

// BuildTeamRosterRoute builds the notification route pointing to a team's roster page.
func BuildTeamRosterRoute(league string, teamID uint) string {
	return fmt.Sprintf("/%s/team/%d", league, teamID)
}

// BuildTeamRecruitingRoute builds the route pointing to a team's recruiting page.
func BuildTeamRecruitingRoute(league string, teamID uint) string {
	return fmt.Sprintf("/%s/recruiting", league)
}

// BuildTeamGameplanRoute builds the route pointing to a team's gameplan page.
func BuildTeamGameplanRoute(league string, teamID uint) string {
	return fmt.Sprintf("/%s/gameplan", league)
}

// BuildTeamAffiliateRosterRoute builds the route pointing to a PHL team's affiliate roster page.
func BuildTeamAffiliateRosterRoute(teamID uint) string {
	return fmt.Sprintf("/phl/team/%d", teamID)
}

// BuildForumThreadRoute builds the route pointing to a specific forum thread.
func BuildForumThreadRoute(threadID string) string {
	return fmt.Sprintf("/forums/thread/%s", threadID)
}

// BuildForumPostRoute builds the route pointing to a specific post inside a thread.
func BuildForumPostRoute(threadID string, postID string) string {
	return fmt.Sprintf("/forums/thread/%s#post-%s", threadID, postID)
}

// BuildSourceEventKey generates a stable idempotency key from its components.
// Example: "injury:chl:season5:game42:player7"
func BuildSourceEventKey(parts ...string) string {
	key := ""
	for i, p := range parts {
		if i > 0 {
			key += ":"
		}
		key += p
	}
	return key
}

// UintToString converts a uint to a decimal string.
func UintToString(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}
