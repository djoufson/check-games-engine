package card_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
)

// TestShouldMatchByColor_WhenPlayingRedJoker tests red joker can be played on red cards
func TestShouldMatchByColor_WhenPlayingRedJoker(t *testing.T) {
	// Arrange
	// Create a Hearts (red) Queen as the top card
	topCard := card.NewCard(card.Hearts, card.Queen)
	redJoker := card.NewRedJoker()

	// Act & Assert
	// Check if red joker matches red card
	if redJoker.Color != topCard.Color {
		t.Errorf("Expected red joker color (%v) to match Hearts card color (%v)", redJoker.Color, topCard.Color)
	}

	if !redJoker.IsWildCard() {
		t.Error("Expected joker to be identified as a wild card")
	}
}

// TestShouldNotMatchByColor_WhenPlayingRedJokerOnBlackCard tests red joker cannot be played on black cards
func TestShouldNotMatchByColor_WhenPlayingRedJokerOnBlackCard(t *testing.T) {
	// Arrange
	// Create a Spades (black) King as the top card
	topCard := card.NewCard(card.Spades, card.King)
	redJoker := card.NewRedJoker()

	// Act & Assert
	// Check if red joker doesn't match black card by color
	if redJoker.Color == topCard.Color {
		t.Errorf("Red joker color (%v) should NOT match Spades card color (%v)", redJoker.Color, topCard.Color)
	}
}

// TestShouldMatchByColor_WhenPlayingBlackJoker tests black joker can be played on black cards
func TestShouldMatchByColor_WhenPlayingBlackJoker(t *testing.T) {
	// Arrange
	// Create a Clubs (black) card as the top card
	topCard := card.NewCard(card.Clubs, card.Ten)
	blackJoker := card.NewBlackJoker()

	// Act & Assert
	// Check if black joker matches black card by color
	if blackJoker.Color != topCard.Color {
		t.Errorf("Expected black joker color (%v) to match Clubs card color (%v)", blackJoker.Color, topCard.Color)
	}

	if !blackJoker.IsWildCard() {
		t.Error("Expected joker to be identified as a wild card")
	}
}

// TestShouldNotMatchByColor_WhenPlayingBlackJokerOnRedCard tests black joker cannot be played on red cards
func TestShouldNotMatchByColor_WhenPlayingBlackJokerOnRedCard(t *testing.T) {
	// Arrange
	// Create a Diamonds (red) card as the top card
	topCard := card.NewCard(card.Diamonds, card.Nine)
	blackJoker := card.NewBlackJoker()

	// Act & Assert
	// Check if black joker doesn't match red card by color
	if blackJoker.Color == topCard.Color {
		t.Errorf("Black joker color (%v) should NOT match Diamonds card color (%v)", blackJoker.Color, topCard.Color)
	}
}

// TestShouldMatchJoker_WhenPlayingJokerOnJoker tests joker can be played on another joker
func TestShouldMatchJoker_WhenPlayingJokerOnJoker(t *testing.T) {
	// Arrange
	// Create both types of jokers
	redJoker := card.NewRedJoker()
	blackJoker := card.NewBlackJoker()

	// Act & Assert
	// Check that both cards are jokers
	if !redJoker.IsJoker() || !blackJoker.IsJoker() {
		t.Error("Both cards should be identified as jokers")
	}

	// The game should allow playing a joker on another joker regardless of color
	// This is a rule validation that would typically be in the state or game package
}
