package card_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
)

// TestCardProperties tests the special properties of different cards
func TestCardProperties(t *testing.T) {
	tests := []struct {
		name          string
		card          card.Card
		isJoker       bool
		isWildCard    bool
		isSkip        bool
		isSuitChanger bool
		isTransparent bool
		drawPenalty   int
	}{
		{
			name:          "Ace of Spades",
			card:          card.NewCard(card.Spades, card.Ace),
			isJoker:       false,
			isWildCard:    false,
			isSkip:        true,
			isSuitChanger: false,
			isTransparent: false,
			drawPenalty:   0,
		},
		{
			name:          "Seven of Hearts",
			card:          card.NewCard(card.Hearts, card.Seven),
			isJoker:       false,
			isWildCard:    true,
			isSkip:        false,
			isSuitChanger: false,
			isTransparent: false,
			drawPenalty:   2,
		},
		{
			name:          "Jack of Diamonds",
			card:          card.NewCard(card.Diamonds, card.Jack),
			isJoker:       false,
			isWildCard:    false,
			isSkip:        false,
			isSuitChanger: true,
			isTransparent: false,
			drawPenalty:   0,
		},
		{
			name:          "Two of Clubs",
			card:          card.NewCard(card.Clubs, card.Two),
			isJoker:       false,
			isWildCard:    false,
			isSkip:        false,
			isSuitChanger: false,
			isTransparent: true,
			drawPenalty:   0,
		},
		{
			name:          "Red Joker",
			card:          card.NewRedJoker(),
			isJoker:       true,
			isWildCard:    true,
			isSkip:        false,
			isSuitChanger: false,
			isTransparent: false,
			drawPenalty:   4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.card.IsJoker() != tt.isJoker {
				t.Errorf("IsJoker() = %v, want %v", tt.card.IsJoker(), tt.isJoker)
			}
			if tt.card.IsWildCard() != tt.isWildCard {
				t.Errorf("IsWildCard() = %v, want %v", tt.card.IsWildCard(), tt.isWildCard)
			}
			if tt.card.IsSkip() != tt.isSkip {
				t.Errorf("IsSkip() = %v, want %v", tt.card.IsSkip(), tt.isSkip)
			}
			if tt.card.IsSuitChanger() != tt.isSuitChanger {
				t.Errorf("IsSuitChanger() = %v, want %v", tt.card.IsSuitChanger(), tt.isSuitChanger)
			}
			if tt.card.IsTransparent() != tt.isTransparent {
				t.Errorf("IsTransparent() = %v, want %v", tt.card.IsTransparent(), tt.isTransparent)
			}
			if tt.card.GetDrawPenalty() != tt.drawPenalty {
				t.Errorf("GetDrawPenalty() = %v, want %v", tt.card.GetDrawPenalty(), tt.drawPenalty)
			}
		})
	}
}
