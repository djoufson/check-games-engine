package player_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/player"
)

// TestPlayerCreation tests creating a new player
func TestPlayerCreation(t *testing.T) {
	p := player.New("player1")
	if p.ID != "player1" {
		t.Errorf("Expected ID to be 'player1', got '%s'", p.ID)
	}
	if len(p.Hand) != 0 {
		t.Errorf("Expected new player to have empty hand, got %d cards", len(p.Hand))
	}
}

// TestPlayerEmptyHandStatus tests the HasEmptyHand method
func TestPlayerEmptyHandStatus(t *testing.T) {
	p := player.New("player1")
	if !p.HasEmptyHand() {
		t.Error("Expected new player to have empty hand")
	}

	p.AddToHand(card.NewCard(card.Spades, card.Ace))
	if p.HasEmptyHand() {
		t.Error("Expected player with cards to not have empty hand")
	}
}
