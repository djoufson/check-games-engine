package game_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/game"
)

// TestShouldCreateGameWithCorrectPlayerCount_WhenUsingDefaultOptions tests creating a game with default options
func TestShouldCreateGameWithCorrectPlayerCount_WhenUsingDefaultOptions(t *testing.T) {
	// Arrange
	playerIDs := []string{"player1", "player2", "player3"}

	// Act
	g, err := game.New(playerIDs, nil)

	// Assert
	if err != nil {
		t.Fatalf("Failed to create new game: %v", err)
	}

	if g.GetPlayerCount() != 3 {
		t.Errorf("Expected 3 players, got %d", g.GetPlayerCount())
	}
}

// TestShouldSetFirstPlayerAsActive_WhenCreatingNewGame tests that the first player is set as active
func TestShouldSetFirstPlayerAsActive_WhenCreatingNewGame(t *testing.T) {
	// Arrange
	playerIDs := []string{"player1", "player2", "player3"}

	// Act
	g, err := game.New(playerIDs, nil)

	// Assert
	if err != nil {
		t.Fatalf("Failed to create new game: %v", err)
	}

	if g.CurrentPlayerID() != "player1" {
		t.Errorf("Expected current player to be player1, got %s", g.CurrentPlayerID())
	}
}

// TestShouldNotBeGameOver_WhenCreatingNewGame tests that a new game is not over
func TestShouldNotBeGameOver_WhenCreatingNewGame(t *testing.T) {
	// Arrange
	playerIDs := []string{"player1", "player2", "player3"}

	// Act
	g, err := game.New(playerIDs, nil)

	// Assert
	if err != nil {
		t.Fatalf("Failed to create new game: %v", err)
	}

	if g.IsGameOver() {
		t.Error("New game should not be over")
	}
}

// TestShouldDistributeCorrectNumberOfCards_WhenUsingCustomOptions tests creating a game with custom options
func TestShouldDistributeCorrectNumberOfCards_WhenUsingCustomOptions(t *testing.T) {
	// Arrange
	playerIDs := []string{"player1", "player2", "player3"}
	opts := &game.Options{
		InitialCards: 5,
		RandomSeed:   12345,
	}

	// Act
	g, err := game.New(playerIDs, opts)

	// Assert
	if err != nil {
		t.Fatalf("Failed to create new game with options: %v", err)
	}

	hand, err := g.GetPlayerHand("player1")
	if err != nil {
		t.Fatalf("Failed to get player hand: %v", err)
	}

	if len(hand) != 5 {
		t.Errorf("Expected hand size to be 5, got %d", len(hand))
	}
}

// TestShouldReturnError_WhenCreatingGameWithTooFewPlayers tests creating a game with too few players
func TestShouldReturnError_WhenCreatingGameWithTooFewPlayers(t *testing.T) {
	// Arrange
	playerIDs := []string{"player1"}

	// Act
	_, err := game.New(playerIDs, nil)

	// Assert
	if err == nil {
		t.Error("Expected error with only one player")
	}
}
