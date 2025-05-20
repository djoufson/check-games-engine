package player_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/player"
)

func TestNew(t *testing.T) {
	p := player.New("player1")
	if p.ID != "player1" {
		t.Errorf("Expected ID to be 'player1', got '%s'", p.ID)
	}
	if len(p.Hand) != 0 {
		t.Errorf("Expected new player to have empty hand, got %d cards", len(p.Hand))
	}
}

func TestAddToHand(t *testing.T) {
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

func TestAddCardsToHand(t *testing.T) {
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

func TestRemoveFromHand(t *testing.T) {
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

func TestHasCard(t *testing.T) {
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

func TestCanPlayCardOn(t *testing.T) {
	tests := []struct {
		name          string
		playedCard    card.Card
		topCard       card.Card
		inAttackChain bool
		expected      bool
	}{
		{
			name:          "Same suit",
			playedCard:    card.NewCard(card.Spades, card.King),
			topCard:       card.NewCard(card.Spades, card.Queen),
			inAttackChain: false,
			expected:      true,
		},
		{
			name:          "Same rank",
			playedCard:    card.NewCard(card.Hearts, card.King),
			topCard:       card.NewCard(card.Spades, card.King),
			inAttackChain: false,
			expected:      true,
		},
		{
			name:          "Different suit and rank",
			playedCard:    card.NewCard(card.Hearts, card.King),
			topCard:       card.NewCard(card.Spades, card.Queen),
			inAttackChain: false,
			expected:      false,
		},
		{
			name:          "Transparent card (2)",
			playedCard:    card.NewCard(card.Hearts, card.Two),
			topCard:       card.NewCard(card.Spades, card.Queen),
			inAttackChain: false,
			expected:      true,
		},
		{
			name:          "Transparent card in attack chain",
			playedCard:    card.NewCard(card.Hearts, card.Two),
			topCard:       card.NewCard(card.Spades, card.Seven),
			inAttackChain: true,
			expected:      false,
		},
		{
			name:          "Wild card on wild card",
			playedCard:    card.NewCard(card.Hearts, card.Seven),
			topCard:       card.NewCard(card.Spades, card.Seven),
			inAttackChain: false,
			expected:      true,
		},
		{
			name:          "Wild card on wild card in attack chain",
			playedCard:    card.NewCard(card.Hearts, card.Seven),
			topCard:       card.NewCard(card.Spades, card.Seven),
			inAttackChain: true,
			expected:      true,
		},
		{
			name:          "Joker on joker",
			playedCard:    card.NewRedJoker(),
			topCard:       card.NewBlackJoker(),
			inAttackChain: false,
			expected:      true,
		},
		{
			name:          "Joker on joker in attack chain",
			playedCard:    card.NewRedJoker(),
			topCard:       card.NewBlackJoker(),
			inAttackChain: true,
			expected:      true,
		},
		{
			name:          "Joker on same color card",
			playedCard:    card.NewRedJoker(),
			topCard:       card.NewCard(card.Hearts, card.King),
			inAttackChain: false,
			expected:      true,
		},
		{
			name:          "Joker on different color card",
			playedCard:    card.NewRedJoker(),
			topCard:       card.NewCard(card.Spades, card.King),
			inAttackChain: false,
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := player.CanPlayCardOn(tt.playedCard, tt.topCard, tt.inAttackChain)
			if result != tt.expected {
				t.Errorf("CanPlayCardOn() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetPlayableCards(t *testing.T) {
	p := player.New("player1")
	// Add some cards to hand
	p.AddCardsToHand([]card.Card{
		card.NewCard(card.Spades, card.Ace),
		card.NewCard(card.Hearts, card.King),
		card.NewCard(card.Diamonds, card.Seven), // Wild card
		card.NewCard(card.Clubs, card.Two),      // Transparent
		card.NewRedJoker(),                      // Wild card
	})

	// Test which cards are playable on Queen of Spades
	topCard := card.NewCard(card.Spades, card.Queen)
	playable := p.GetPlayableCards(topCard, false)

	// Playable cards should be:
	// 1. Ace of Spades (same suit)
	// 2. Two of Clubs (transparent)
	// 3. Seven of Diamonds and Red Joker (if we consider that wild cards can be played anytime)
	if len(playable) != 2 {
		t.Errorf("Expected 2 playable cards, got %d", len(playable))
	}

	// Test during attack chain
	topCard = card.NewCard(card.Diamonds, card.Seven) // Wild Card
	playable = p.GetPlayableCards(topCard, true)

	// Playable cards should be:
	// 1. Seven of Diamonds
	// 2. Red Joker
	if len(playable) != 2 {
		t.Errorf("Expected 2 playable cards during attack chain, got %d", len(playable))
	}
}

func TestHasEmptyHand(t *testing.T) {
	p := player.New("player1")
	if !p.HasEmptyHand() {
		t.Error("Expected new player to have empty hand")
	}

	p.AddToHand(card.NewCard(card.Spades, card.Ace))
	if p.HasEmptyHand() {
		t.Error("Expected player with cards to not have empty hand")
	}
}

func TestHasMatchingCard(t *testing.T) {
	p := player.New("player1")
	p.AddCardsToHand([]card.Card{
		card.NewCard(card.Spades, card.Ace),
		card.NewCard(card.Hearts, card.King),
	})

	// Should match Ace of Spades with Queen of Spades (same suit)
	if !p.HasMatchingCard(card.NewCard(card.Spades, card.Queen), false) {
		t.Error("Expected player to have a matching card (same suit)")
	}

	// Should match King of Hearts with King of Clubs (same rank)
	if !p.HasMatchingCard(card.NewCard(card.Clubs, card.King), false) {
		t.Error("Expected player to have a matching card (same rank)")
	}

	// Should not match with a card with different suit and rank
	if p.HasMatchingCard(card.NewCard(card.Diamonds, card.Queen), false) {
		t.Error("Expected player not to have a matching card")
	}
}
