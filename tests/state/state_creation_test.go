package state_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/state"
)

// TestShouldCreateStateWithCorrectPlayerCount_WhenUsingDefaultOptions tests player count with default options
func TestShouldCreateStateWithCorrectPlayerCount_WhenUsingDefaultOptions(t *testing.T) {
	// Arrange
	playerIDs := []string{"player1", "player2", "player3"}

	// Act
	gameState, err := state.New(playerIDs, nil)

	// Assert
	if err != nil {
		t.Fatalf("Failed to create new game: %v", err)
	}

	if len(gameState.ActivePlayers) != 3 {
		t.Errorf("Expected 3 players, got %d", len(gameState.ActivePlayers))
	}
}

// TestShouldSetFirstPlayerAsActive_WhenCreatingNewState tests that the first player is active
func TestShouldSetFirstPlayerAsActive_WhenCreatingNewState(t *testing.T) {
	// Arrange
	playerIDs := []string{"player1", "player2", "player3"}

	// Act
	gameState, err := state.New(playerIDs, nil)

	// Assert
	if err != nil {
		t.Fatalf("Failed to create new game: %v", err)
	}

	if gameState.CurrentPlayerID() != "player1" {
		t.Errorf("Expected current player to be player1, got %s", gameState.CurrentPlayerID())
	}
}

// TestShouldNotBeGameOver_WhenCreatingNewState tests that a new game is not over
func TestShouldNotBeGameOver_WhenCreatingNewState(t *testing.T) {
	// Arrange
	playerIDs := []string{"player1", "player2", "player3"}

	// Act
	gameState, err := state.New(playerIDs, nil)

	// Assert
	if err != nil {
		t.Fatalf("Failed to create new game: %v", err)
	}

	if gameState.IsGameOver() {
		t.Error("New game should not be over")
	}
}

// TestShouldDistributeCorrectNumberOfCards_WhenUsingCustomOptions tests card distribution with custom options
func TestShouldDistributeCorrectNumberOfCards_WhenUsingCustomOptions(t *testing.T) {
	// Arrange
	playerIDs := []string{"player1", "player2", "player3"}
	opts := &state.GameOptions{
		InitialCards: 5,
		RandomSeed:   12345,
	}

	// Act
	gameState, err := state.New(playerIDs, opts)

	// Assert
	if err != nil {
		t.Fatalf("Failed to create new game with options: %v", err)
	}

	player1 := gameState.FindPlayerByID("player1")
	if player1 == nil {
		t.Fatalf("Failed to find player1")
	}

	if len(player1.Hand) != 5 {
		t.Errorf("Expected hand size to be 5, got %d", len(player1.Hand))
	}
}

// TestShouldReturnError_WhenCreatingStateWithTooFewPlayers tests too few players error
func TestShouldReturnError_WhenCreatingStateWithTooFewPlayers(t *testing.T) {
	// Arrange
	playerIDs := []string{"player1"}

	// Act
	_, err := state.New(playerIDs, nil)

	// Assert
	if err == nil {
		t.Error("Expected error with only one player")
	}
}
