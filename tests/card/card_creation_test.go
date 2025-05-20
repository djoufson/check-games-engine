package card_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
)

// TestNewCard_ShouldSetSuitCorrectly_WhenCreatingCard verifies that the suit is set correctly
func TestNewCard_ShouldSetSuitCorrectly_WhenCreatingCard(t *testing.T) {
	// Arrange & Act
	c := card.NewCard(card.Spades, card.Ace)

	// Assert
	if c.Suit != card.Spades {
		t.Errorf("Expected suit to be Spades, got %v", c.Suit)
	}
}

// TestNewCard_ShouldSetRankCorrectly_WhenCreatingCard verifies that the rank is set correctly
func TestNewCard_ShouldSetRankCorrectly_WhenCreatingCard(t *testing.T) {
	// Arrange & Act
	c := card.NewCard(card.Spades, card.Ace)

	// Assert
	if c.Rank != card.Ace {
		t.Errorf("Expected rank to be Ace, got %v", c.Rank)
	}
}

// TestNewCard_ShouldSetBlackColor_WhenSuitIsSpades verifies color is automatically set to black for spades
func TestNewCard_ShouldSetBlackColor_WhenSuitIsSpades(t *testing.T) {
	// Arrange & Act
	c := card.NewCard(card.Spades, card.Ace)

	// Assert
	if c.Color != card.Black {
		t.Errorf("Expected color to be Black, got %v", c.Color)
	}
}

// TestNewCard_ShouldSetRedColor_WhenSuitIsHearts verifies color is automatically set to red for hearts
func TestNewCard_ShouldSetRedColor_WhenSuitIsHearts(t *testing.T) {
	// Arrange & Act
	c := card.NewCard(card.Hearts, card.King)

	// Assert
	if c.Color != card.Red {
		t.Errorf("Expected color to be Red, got %v", c.Color)
	}
}
