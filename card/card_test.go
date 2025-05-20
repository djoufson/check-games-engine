package card

import (
	"testing"
)

func TestNewCard(t *testing.T) {
	// Test regular cards
	card := NewCard(Spades, Ace)
	if card.Suit != Spades {
		t.Errorf("Expected suit to be Spades, got %v", card.Suit)
	}
	if card.Rank != Ace {
		t.Errorf("Expected rank to be Ace, got %v", card.Rank)
	}
	if card.Color != Black {
		t.Errorf("Expected color to be Black, got %v", card.Color)
	}

	// Test hearts card (should be red)
	card = NewCard(Hearts, King)
	if card.Color != Red {
		t.Errorf("Expected color to be Red, got %v", card.Color)
	}
}

func TestJokers(t *testing.T) {
	redJoker := NewRedJoker()
	if redJoker.Suit != Joker || redJoker.Color != Red {
		t.Errorf("Expected Red Joker, got %v", redJoker)
	}

	blackJoker := NewBlackJoker()
	if blackJoker.Suit != Joker || blackJoker.Color != Black {
		t.Errorf("Expected Black Joker, got %v", blackJoker)
	}
}

func TestStringRepresentation(t *testing.T) {
	ace := NewCard(Spades, Ace)
	if ace.String() != "ACE of SPADES" {
		t.Errorf("Expected 'ACE of SPADES', got '%s'", ace.String())
	}

	joker := NewRedJoker()
	if joker.String() != "RED Joker" {
		t.Errorf("Expected 'RED Joker', got '%s'", joker.String())
	}
}

func TestSpecialCardChecks(t *testing.T) {
	tests := []struct {
		name       string
		card       Card
		isJoker    bool
		isWildCard bool
		isSkip     bool
		isSuitChanger bool
		isTransparent bool
		drawPenalty int
	}{
		{
			name:       "Ace of Spades",
			card:       NewCard(Spades, Ace),
			isJoker:    false,
			isWildCard: false,
			isSkip:     true,
			isSuitChanger: false,
			isTransparent: false,
			drawPenalty: 0,
		},
		{
			name:       "Seven of Hearts",
			card:       NewCard(Hearts, Seven),
			isJoker:    false,
			isWildCard: true,
			isSkip:     false,
			isSuitChanger: false,
			isTransparent: false,
			drawPenalty: 2,
		},
		{
			name:       "Jack of Diamonds",
			card:       NewCard(Diamonds, Jack),
			isJoker:    false,
			isWildCard: false,
			isSkip:     false,
			isSuitChanger: true,
			isTransparent: false,
			drawPenalty: 0,
		},
		{
			name:       "Two of Clubs",
			card:       NewCard(Clubs, Two),
			isJoker:    false,
			isWildCard: false,
			isSkip:     false,
			isSuitChanger: false,
			isTransparent: true,
			drawPenalty: 0,
		},
		{
			name:       "Red Joker",
			card:       NewRedJoker(),
			isJoker:    true,
			isWildCard: true,
			isSkip:     false,
			isSuitChanger: false,
			isTransparent: false,
			drawPenalty: 4,
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