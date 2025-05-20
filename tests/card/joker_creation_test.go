package card_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
)

// TestJokerCreation tests the creation and properties of joker cards
func TestJokerCreation(t *testing.T) {
	redJoker := card.NewRedJoker()
	if redJoker.Suit != card.Joker || redJoker.Color != card.Red {
		t.Errorf("Expected Red Joker, got %v", redJoker)
	}

	blackJoker := card.NewBlackJoker()
	if blackJoker.Suit != card.Joker || blackJoker.Color != card.Black {
		t.Errorf("Expected Black Joker, got %v", blackJoker)
	}
}
