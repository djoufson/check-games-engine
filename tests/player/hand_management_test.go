package player_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/player"
)

// TestAddCardToHand tests adding a single card to a player's hand
func TestAddCardToHand(t *testing.T) {
	p := player.New("player1")
	c := card.NewCard(card.Spades, card.Ace)

	p.AddToHand(c)
	if len(p.Hand) != 1 {
		t.Errorf("Expected hand to have 1 card, got %d", len(p.Hand))
	}
	if p.Hand[0].Suit != card.Spades || p.Hand[0].Rank != card.Ace {
		t.Errorf("Expected Ace of Spades, got %v", p.Hand[0])
	}
}

// TestAddMultipleCardsToHand tests adding multiple cards to a player's hand
func TestAddMultipleCardsToHand(t *testing.T) {
	p := player.New("player1")
	cards := []card.Card{
		card.NewCard(card.Spades, card.Ace),
		card.NewCard(card.Hearts, card.King),
	}

	p.AddCardsToHand(cards)
	if len(p.Hand) != 2 {
		t.Errorf("Expected hand to have 2 cards, got %d", len(p.Hand))
	}
}

// TestRemoveCardFromHand tests removing a card from a player's hand
func TestRemoveCardFromHand(t *testing.T) {
	p := player.New("player1")
	aceSpades := card.NewCard(card.Spades, card.Ace)
	kingHearts := card.NewCard(card.Hearts, card.King)

	p.AddCardsToHand([]card.Card{aceSpades, kingHearts})

	// Remove a card that exists in the hand
	removedCard, ok := p.RemoveFromHand(aceSpades)
	if !ok {
		t.Error("Expected card to be removed successfully")
	}
	if removedCard.Suit != card.Spades || removedCard.Rank != card.Ace {
		t.Errorf("Expected to remove Ace of Spades, got %v", removedCard)
	}
	if len(p.Hand) != 1 {
		t.Errorf("Expected hand to have 1 card after removal, got %d", len(p.Hand))
	}

	// Try to remove a card that doesn't exist in the hand
	_, ok = p.RemoveFromHand(aceSpades)
	if ok {
		t.Error("Expected removal of non-existent card to fail")
	}
}

// TestCheckCardInHand tests checking if a card is in a player's hand
func TestCheckCardInHand(t *testing.T) {
	p := player.New("player1")
	aceSpades := card.NewCard(card.Spades, card.Ace)
	kingHearts := card.NewCard(card.Hearts, card.King)

	p.AddCardsToHand([]card.Card{aceSpades})

	if !p.HasCard(aceSpades) {
		t.Error("Expected player to have Ace of Spades")
	}

	if p.HasCard(kingHearts) {
		t.Error("Expected player not to have King of Hearts")
	}
}
