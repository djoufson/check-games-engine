package state_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/deck"
	"github.com/djoufson/check-games-engine/player"
	"github.com/djoufson/check-games-engine/state"
)

func TestAttackChain(t *testing.T) {
	// Create players with specific cards
	player1 := player.New("player1")
	player1.AddCardsToHand([]card.Card{
		card.NewCard(card.Hearts, card.Seven), // Wild (+2)
		card.NewRedJoker(),                    // Wild (+4)
	})

	player2 := player.New("player2")
	player2.AddCardsToHand([]card.Card{
		card.NewCard(card.Spades, card.Seven), // Wild (+2)
		card.NewCard(card.Clubs, card.King),   // Regular card
	})

	// Create a state with known top card
	drawPile := deck.New()
	topCard := card.NewCard(card.Hearts, card.Queen)

	gameState := &state.State{
		Players:            []*player.Player{player1, player2},
		ActivePlayers:      []string{"player1", "player2"},
		CurrentPlayerIndex: 0,
		Direction:          state.Clockwise,
		DrawPile:           drawPile,
		DiscardPile:        []card.Card{topCard},
		TopCard:            topCard,
		InAttackChain:      false,
		AttackAmount:       0,
		LastActiveSuit:     card.Hearts,
	}

	// Play a wild card to start an attack chain
	err := gameState.PlayCard("player1", card.NewCard(card.Hearts, card.Seven))
	if err != nil {
		t.Fatalf("Failed to play Seven: %v", err)
	}

	// Check attack chain started
	if !gameState.InAttackChain {
		t.Error("Expected to be in attack chain")
	}

	if gameState.AttackAmount != 2 {
		t.Errorf("Expected attack amount to be 2, got %d", gameState.AttackAmount)
	}

	// Check that it's player2's turn now
	if gameState.CurrentPlayerID() != "player2" {
		t.Errorf("Expected current player to be player2, got %s", gameState.CurrentPlayerID())
	}

	// Player2 responds with another wild card
	err = gameState.PlayCard("player2", card.NewCard(card.Spades, card.Seven))
	if err != nil {
		t.Fatalf("Failed to play Seven to continue attack chain: %v", err)
	}

	// Check attack amount increased
	if gameState.AttackAmount != 4 {
		t.Errorf("Expected attack amount to be 4, got %d", gameState.AttackAmount)
	}

	// Check that it's player1's turn now
	if gameState.CurrentPlayerID() != "player1" {
		t.Errorf("Expected current player to be player1, got %s", gameState.CurrentPlayerID())
	}

	// Player1 plays Joker to continue attack chain
	err = gameState.PlayCard("player1", card.NewRedJoker())
	if err != nil {
		t.Fatalf("Failed to play Joker: %v", err)
	}

	// Check attack amount increased
	if gameState.AttackAmount != 8 {
		t.Errorf("Expected attack amount to be 8, got %d", gameState.AttackAmount)
	}

	// Get initial hand size before player2 draws
	initialHandSize := len(player2.Hand)

	// Player2 has no more wild cards, so must draw
	err = gameState.DrawCard("player2")
	if err != nil {
		t.Fatalf("Failed to draw cards: %v", err)
	}

	// Check player2 drew 8 cards
	if len(player2.Hand) != initialHandSize+8 {
		t.Errorf("Expected hand size to increase by 8, got %d", len(player2.Hand))
	}

	// Check attack chain ended
	if gameState.InAttackChain {
		t.Error("Expected attack chain to end")
	}
}
