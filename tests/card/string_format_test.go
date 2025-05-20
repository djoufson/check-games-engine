package card_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
)

// TestShouldReturnFormattedString_WhenConvertingRegularCardToString tests string format for regular cards
func TestShouldReturnFormattedString_WhenConvertingRegularCardToString(t *testing.T) {
	// Arrange
	ace := card.NewCard(card.Spades, card.Ace)

	// Act
	result := ace.String()

	// Assert
	if result != "ACE of SPADES" {
		t.Errorf("Expected 'ACE of SPADES', got '%s'", result)
	}
}

// TestShouldReturnFormattedString_WhenConvertingJokerToString tests string format for joker cards
func TestShouldReturnFormattedString_WhenConvertingJokerToString(t *testing.T) {
	// Arrange
	joker := card.NewRedJoker()

	// Act
	result := joker.String()

	// Assert
	if result != "RED Joker" {
		t.Errorf("Expected 'RED Joker', got '%s'", result)
	}
}
