package card_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
)

func TestNewCard(t *testing.T) {
	// Test regular cards
	c := card.NewCard(card.Spades, card.Ace)
	if c.Suit != card.Spades {
		t.Errorf("Expected suit to be Spades, got %v", c.Suit)
	}
	if c.Rank != card.Ace {
		t.Errorf("Expected rank to be Ace, got %v", c.Rank)
	}
	if c.Color != card.Black {
		t.Errorf("Expected color to be Black, got %v", c.Color)
	}

	// Test hearts card (should be red)
	c = card.NewCard(card.Hearts, card.King)
	if c.Color != card.Red {
		t.Errorf("Expected color to be Red, got %v", c.Color)
	}
}
