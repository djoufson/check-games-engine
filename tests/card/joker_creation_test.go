package card_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
)

// TestShouldCreateRedJokerWithCorrectProperties_WhenCallingNewRedJoker tests red joker creation
func TestShouldCreateRedJokerWithCorrectProperties_WhenCallingNewRedJoker(t *testing.T) {
	// Arrange & Act
	redJoker := card.NewRedJoker()

	// Assert
	if redJoker.Suit != card.Joker {
		t.Errorf("Expected Joker suit, got %v", redJoker.Suit)
	}

	if redJoker.Color != card.Red {
		t.Errorf("Expected Red color, got %v", redJoker.Color)
	}
}

// TestShouldCreateBlackJokerWithCorrectProperties_WhenCallingNewBlackJoker tests black joker creation
func TestShouldCreateBlackJokerWithCorrectProperties_WhenCallingNewBlackJoker(t *testing.T) {
	// Arrange & Act
	blackJoker := card.NewBlackJoker()

	// Assert
	if blackJoker.Suit != card.Joker {
		t.Errorf("Expected Joker suit, got %v", blackJoker.Suit)
	}

	if blackJoker.Color != card.Black {
		t.Errorf("Expected Black color, got %v", blackJoker.Color)
	}
}
