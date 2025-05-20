package game_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/game"
)

func TestNewGame(t *testing.T) {
	playerIDs := []string{"player1", "player2", "player3"}

	// Test with default options
	g, err := game.New(playerIDs, nil)
	if err != nil {
		t.Fatalf("Failed to create new game: %v", err)
	}

	// Check basic game state
	if g.GetPlayerCount() != 3 {
		t.Errorf("Expected 3 players, got %d", g.GetPlayerCount())
	}

	if g.CurrentPlayerID() != "player1" {
		t.Errorf("Expected current player to be player1, got %s", g.CurrentPlayerID())
	}

	if g.IsGameOver() {
		t.Error("New game should not be over")
	}

	// Test with custom options
	opts := &game.Options{
		InitialCards: 5,
		RandomSeed:   12345,
	}

	g, err = game.New(playerIDs, opts)
	if err != nil {
		t.Fatalf("Failed to create new game with options: %v", err)
	}

	// Check that player hands have expected size
	hand, err := g.GetPlayerHand("player1")
	if err != nil {
		t.Fatalf("Failed to get player hand: %v", err)
	}

	if len(hand) != 5 {
		t.Errorf("Expected hand size to be 5, got %d", len(hand))
	}

	// Test with invalid player count
	_, err = game.New([]string{"player1"}, nil)
	if err == nil {
		t.Error("Expected error with only one player")
	}
}
