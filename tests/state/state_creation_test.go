package state_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/state"
)

func TestStateCreation(t *testing.T) {
	playerIDs := []string{"player1", "player2", "player3"}

	// Test with default options
	gameState, err := state.New(playerIDs, nil)
	if err != nil {
		t.Fatalf("Failed to create new game: %v", err)
	}

	// Check basic game state
	if len(gameState.ActivePlayers) != 3 {
		t.Errorf("Expected 3 players, got %d", len(gameState.ActivePlayers))
	}

	if gameState.CurrentPlayerID() != "player1" {
		t.Errorf("Expected current player to be player1, got %s", gameState.CurrentPlayerID())
	}

	if gameState.IsGameOver() {
		t.Error("New game should not be over")
	}

	// Test with custom options
	opts := &state.GameOptions{
		InitialCards: 5,
		RandomSeed:   12345,
	}

	gameState, err = state.New(playerIDs, opts)
	if err != nil {
		t.Fatalf("Failed to create new game with options: %v", err)
	}

	// Check that player hands have expected size
	player1 := gameState.FindPlayerByID("player1")
	if player1 == nil {
		t.Fatalf("Failed to find player1")
	}

	if len(player1.Hand) != 5 {
		t.Errorf("Expected hand size to be 5, got %d", len(player1.Hand))
	}

	// Test with invalid player count
	_, err = state.New([]string{"player1"}, nil)
	if err == nil {
		t.Error("Expected error with only one player")
	}
}
