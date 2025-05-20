package player_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/player"
)

// TestShouldSetCorrectID_WhenCreatingNewPlayer tests that the player ID is set correctly
func TestShouldSetCorrectID_WhenCreatingNewPlayer(t *testing.T) {
	// Arrange & Act
	p := player.New("player1")

	// Assert
	if p.ID != "player1" {
		t.Errorf("Expected ID to be 'player1', got '%s'", p.ID)
	}
}

// TestShouldHaveEmptyHand_WhenCreatingNewPlayer tests that a new player starts with no cards
func TestShouldHaveEmptyHand_WhenCreatingNewPlayer(t *testing.T) {
	// Arrange & Act
	p := player.New("player1")

	// Assert
	if len(p.Hand) != 0 {
		t.Errorf("Expected new player to have empty hand, got %d cards", len(p.Hand))
	}
}

// TestShouldReportEmptyHand_WhenPlayerHasNoCards tests the HasEmptyHand method with empty hand
func TestShouldReportEmptyHand_WhenPlayerHasNoCards(t *testing.T) {
	// Arrange
	p := player.New("player1")

	// Act & Assert
	if !p.HasEmptyHand() {
		t.Error("Expected new player to have empty hand")
	}
}

// TestShouldNotReportEmptyHand_WhenPlayerHasCards tests the HasEmptyHand method with cards
func TestShouldNotReportEmptyHand_WhenPlayerHasCards(t *testing.T) {
	// Arrange
	p := player.New("player1")
	p.AddToHand(card.NewCard(card.Spades, card.Ace))

	// Act & Assert
	if p.HasEmptyHand() {
		t.Error("Expected player with cards to not have empty hand")
	}
}
