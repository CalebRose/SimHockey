package firebase

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// Collection name constants for the live scoreboard (games only — plays are
// served from the API, not stored in Firebase).
const (
	CollectionCHLGames = "live_chl_games"
	CollectionPHLGames = "live_phl_games"
)

// liveGamesCollection returns the games collection name for the given league.
// league must be "chl" or "phl".
func liveGamesCollection(league string) string {
	if league == "chl" {
		return CollectionCHLGames
	}
	return CollectionPHLGames
}

// PurgeStaleLiveGames deletes all documents from the live games collection where
// IsRevealed == true.  Called at the start of RunGames to clear already-broadcast
// records so the scoreboard only shows fresh, unrevealed games.
func PurgeStaleLiveGames(ctx context.Context, league string) error {
	client := GetFirestoreClient()
	gamesCol := liveGamesCollection(league)

	iter := client.Collection(gamesCol).Where("IsRevealed", "==", true).Documents(ctx)
	defer iter.Stop()

	batch := client.Batch()
	count := 0
	for {
		docSnap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("firebase: PurgeStaleLiveGames(%s) iterate: %w", league, err)
		}
		batch.Delete(docSnap.Ref)
		count++
	}

	if count == 0 {
		return nil
	}

	if _, err := batch.Commit(ctx); err != nil {
		return fmt.Errorf("firebase: PurgeStaleLiveGames(%s) commit: %w", league, err)
	}
	log.Printf("firebase: purged %d stale live game records for league=%s", count, league)
	return nil
}

// UploadLiveGame writes a single game metadata record to the live games
// collection (document ID == GameID).  Called by the StreamScheduler each time
// a game is promoted into an active streaming slot.
func UploadLiveGame(ctx context.Context, game LiveGameRecord, league string) error {
	client := GetFirestoreClient()
	gamesCol := liveGamesCollection(league)
	docID := strconv.FormatUint(uint64(game.GameID), 10)
	if _, err := client.Collection(gamesCol).Doc(docID).Set(ctx, game); err != nil {
		return fmt.Errorf("firebase: UploadLiveGame(gameID=%d, league=%s): %w", game.GameID, league, err)
	}
	return nil
}

// SetGameRevealed marks a live game document as IsRevealed = true in Firestore.
// Called by the StreamScheduler when a game's StreamEndTime has passed.
func SetGameRevealed(ctx context.Context, gameID uint, league string) error {
	client := GetFirestoreClient()
	gamesCol := liveGamesCollection(league)

	docID := strconv.FormatUint(uint64(gameID), 10)
	docRef := client.Collection(gamesCol).Doc(docID)

	_, err := docRef.Update(ctx, []firestore.Update{
		{Path: "IsRevealed", Value: true},
	})
	if err != nil {
		return fmt.Errorf("firebase: SetGameRevealed(gameID=%d, league=%s): %w", gameID, league, err)
	}
	return nil
}

