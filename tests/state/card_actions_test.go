package state_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/deck"
	"github.com/djoufson/check-games-engine/player"
	"github.com/djoufson/check-games-engine/state"
)

// TestCardPlay tests the ability to play cards and have them correctly affect the game state
func TestCardPlay(t *testing.T) {
	// Create players with specific cards
	player1 := player.New("player1")
	player1.AddCardsToHand([]card.Card{
		card.NewCard(card.Spades, card.Ace),
		card.NewCard(card.Hearts, card.King),
	})

	player2 := player.New("player2")
	player2.AddCardsToHand([]card.Card{
		card.NewCard(card.Diamonds, card.Queen),
		card.NewCard(card.Clubs, card.Jack),
	})

	// Create a state with known top card
	gameState := &state.State{
		Players:            []*player.Player{player1, player2},
		ActivePlayers:      []string{"player1", "player2"},
		CurrentPlayerIndex: 0,
		Direction:          state.Clockwise,
		DrawPile:           nil, // We won't be drawing in this test
		DiscardPile:        []card.Card{card.NewCard(card.Spades, card.Queen)},
		TopCard:            card.NewCard(card.Spades, card.Queen),
		InAttackChain:      false,
		AttackAmount:       0,
		LastActiveSuit:     card.Spades,
	}

	// Play a card that matches suit with the top card
	err := gameState.PlayCard("player1", card.NewCard(card.Spades, card.Ace))
	if err != nil {
		t.Fatalf("Failed to play matching card: %v", err)
	}

	// Check that the top card changed
	if gameState.TopCard.Suit != card.Spades || gameState.TopCard.Rank != card.Ace {
		t.Errorf("Expected top card to be Ace of Spades, got %v", gameState.TopCard)
	}

	// Check that it's player2's turn now (Ace skips, but with 2 players it goes back to player1)
	// But since Ace skips the next player, and there are only 2 players, it should be player1's turn again
	if gameState.CurrentPlayerID() != "player1" {
		t.Errorf("Expected current player to be player1 (after Ace skip), got %s", gameState.CurrentPlayerID())
	}

	// Check that player1's hand no longer contains the Ace
	if player1.HasCard(card.NewCard(card.Spades, card.Ace)) {
		t.Error("Expected player1's hand to no longer contain Ace of Spades")
	}
}

// TestCardDrawing tests drawing cards from the deck into a player's hand
func TestCardDrawing(t *testing.T) {
	// Set up a game state for testing DrawCard
	player1 := player.New("player1")
	player2 := player.New("player2")

	// Create a new deck with default cards
	drawPile := deck.New()

	// Get the top card for the initial discard pile
	topCard, _ := drawPile.Draw()

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
		LastActiveSuit:     topCard.Suit,
	}

	// Get initial hand size
	initialHandSize := len(player1.Hand)

	// Draw a card
	err := gameState.DrawCard("player1")
	if err != nil {
		t.Fatalf("Failed to draw card: %v", err)
	}

	// Check that player1's hand has one more card
	if len(player1.Hand) != initialHandSize+1 {
		t.Errorf("Expected hand size to increase by 1, got %d", len(player1.Hand))
	}

	// Check that it's player2's turn now
	if gameState.CurrentPlayerID() != "player2" {
		t.Errorf("Expected current player to be player2, got %s", gameState.CurrentPlayerID())
	}
}
