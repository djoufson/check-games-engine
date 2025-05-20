package game_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/deck"
	"github.com/djoufson/check-games-engine/game"
	"github.com/djoufson/check-games-engine/player"
	"github.com/djoufson/check-games-engine/state"
)

func TestSpecialCardEffects(t *testing.T) {
	// Set up players with specific cards for testing
	jackCard := card.NewCard(card.Spades, card.Jack)   // Suit change - same suit as top card
	aceCard := card.NewCard(card.Spades, card.Ace)     // Skip
	sevenCard := card.NewCard(card.Hearts, card.Seven) // Wild (+2)
	jokerCard := card.NewRedJoker()                    // Wild (+4)

	player1 := player.New("player1")
	player1.AddCardsToHand([]card.Card{
		aceCard,
		sevenCard,
		jackCard,
		jokerCard,
	})

	player2 := player.New("player2")
	player2.AddCardsToHand([]card.Card{
		card.NewCard(card.Hearts, card.Five),
		card.NewCard(card.Spades, card.Seven), // Wild (+2)
		card.NewCard(card.Clubs, card.Jack),   // Suit change
		card.NewBlackJoker(),                  // Wild (+4)
	})

	// Create a state with these players and specific top card
	stateObj := &state.State{
		Players:            []*player.Player{player1, player2},
		ActivePlayers:      []string{player1.ID, player2.ID},
		CurrentPlayerIndex: 0,
		Direction:          state.Clockwise,
		DrawPile:           deck.New(), // Standard deck
		DiscardPile:        []card.Card{card.NewCard(card.Spades, card.Queen)},
		TopCard:            card.NewCard(card.Spades, card.Queen),
		InAttackChain:      false,
		AttackAmount:       0,
		LastActiveSuit:     card.Spades,
	}

	// Initialize a game from this state
	g := game.FromState(stateObj)

	// Test Ace (Skip)
	err := g.PlayCard(player1.ID, aceCard)
	if err != nil {
		t.Fatalf("Failed to play Ace: %v", err)
	}

	// With 2 players, it should stay player1's turn
	if g.CurrentPlayerID() != "player1" {
		t.Errorf("Expected current player to be player1 after Ace, got %s", g.CurrentPlayerID())
	}

	// Test Jack (Suit change)
	// First play Jack
	err = g.PlayCard(player1.ID, jackCard)
	if err != nil {
		t.Fatalf("Failed to play Jack: %v", err)
	}

	// Then change suit
	err = g.ChangeSuit(player1.ID, card.Hearts)
	if err != nil {
		t.Fatalf("Failed to change suit: %v", err)
	}

	// LastActiveSuit should be Hearts now
	if g.GetLastActiveSuit() != card.Hearts {
		t.Errorf("Expected last active suit to be Hearts, got %v", g.GetLastActiveSuit())
	}

	// Test wild card (Seven) - should be player2's turn now
	player2SevenCard := card.NewCard(card.Spades, card.Seven)
	err = g.PlayCard(player2.ID, player2SevenCard)
	if err != nil {
		t.Fatalf("Failed to play Seven: %v", err)
	}

	// Should be in attack chain
	if !g.IsInAttackChain() {
		t.Error("Expected to be in attack chain")
	}

	if g.GetAttackAmount() != 2 {
		t.Errorf("Expected attack amount to be 2, got %d", g.GetAttackAmount())
	}

	// Test wild card (Joker) in attack chain
	err = g.PlayCard(player1.ID, jokerCard)
	if err != nil {
		t.Fatalf("Failed to play Joker: %v", err)
	}

	// Attack amount should be increased
	if g.GetAttackAmount() != 6 {
		t.Errorf("Expected attack amount to be 6, got %d", g.GetAttackAmount())
	}
}
