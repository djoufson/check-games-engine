package card_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
)

// TestCardStringFormat tests the string representation of cards
func TestCardStringFormat(t *testing.T) {
	ace := card.NewCard(card.Spades, card.Ace)
	if ace.String() != "ACE of SPADES" {
		t.Errorf("Expected 'ACE of SPADES', got '%s'", ace.String())
	}

	joker := card.NewRedJoker()
	if joker.String() != "RED Joker" {
		t.Errorf("Expected 'RED Joker', got '%s'", joker.String())
	}
}
