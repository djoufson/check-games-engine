package game_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/game"
)

func TestJSONSerialization(t *testing.T) {
	// Create a game
	playerIDs := []string{"player1", "player2"}
	opts := &game.Options{
		InitialCards: 7,
		RandomSeed:   12345,
	}

	g, _ := game.New(playerIDs, opts)

	// Serialize to JSON
	data, err := g.ToJSON()
	if err != nil {
		t.Fatalf("Failed to serialize game: %v", err)
	}

	// Create a new game from the JSON
	g2, err := game.FromJSON(data)
	if err != nil {
		t.Fatalf("Failed to deserialize game: %v", err)
	}

	// Check that the games match
	if g.CurrentPlayerID() != g2.CurrentPlayerID() {
		t.Errorf("Current player mismatch: %s vs %s",
			g.CurrentPlayerID(), g2.CurrentPlayerID())
	}

	if g.GetPlayerCount() != g2.GetPlayerCount() {
		t.Errorf("Player count mismatch: %d vs %d",
			g.GetPlayerCount(), g2.GetPlayerCount())
	}

	// Check that the top cards match
	topCard1 := g.GetTopCard()
	topCard2 := g2.GetTopCard()

	if topCard1.Suit != topCard2.Suit || topCard1.Rank != topCard2.Rank {
		t.Errorf("Top card mismatch: %v vs %v", topCard1, topCard2)
	}
}
