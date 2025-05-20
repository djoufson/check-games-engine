package state_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/deck"
	"github.com/djoufson/check-games-engine/player"
	"github.com/djoufson/check-games-engine/state"
)

// TestShouldReshuffleDiscardPile_WhenDrawPileIsEmpty tests reshuffling the discard pile
func TestShouldReshuffleDiscardPile_WhenDrawPileIsEmpty(t *testing.T) {
	// Arrange
	player1 := player.New("player1")

	// Create a state with empty draw pile and multiple cards in discard pile
	discardPile := []card.Card{
		card.NewCard(card.Hearts, card.Five),
		card.NewCard(card.Diamonds, card.Six),
		card.NewCard(card.Clubs, card.Seven),
		card.NewCard(card.Spades, card.Queen), // This will be the top card
	}

	topCard := discardPile[len(discardPile)-1] // Queen of Spades

	// Create an empty draw pile (no cards)
	emptyDrawPile := deck.New()
	// Remove all cards from the draw pile
	for !emptyDrawPile.IsEmpty() {
		emptyDrawPile.Draw()
	}

	gameState := &state.State{
		Players:         []*player.Player{player1},
		ActivePlayers:   []string{player1.ID},
		CurrentPlayerId: player1.ID,
		Direction:       state.Clockwise,
		DrawPile:        emptyDrawPile, // Empty draw pile
		DiscardPile:     discardPile,
		TopCard:         topCard,
		InAttackChain:   false,
		AttackAmount:    0,
		LastActiveSuit:  card.Spades,
	}

	// Act
	err := gameState.DrawCard("player1")

	// Assert
	if err != nil {
		t.Fatalf("Failed to draw card after reshuffling: %v", err)
	}

	// Check that draw pile was refilled
	if gameState.DrawPile.IsEmpty() {
		t.Error("Expected draw pile to be refilled")
	}

	// Check that discard pile only contains top card
	if len(gameState.DiscardPile) != 1 {
		t.Errorf("Expected discard pile to only contain top card, got %d cards", len(gameState.DiscardPile))
	}

	// Check that top card remains unchanged
	if gameState.TopCard.Suit != card.Spades || gameState.TopCard.Rank != card.Queen {
		t.Errorf("Expected top card to remain Queen of Spades, got %v", gameState.TopCard)
	}
}

// TestShouldHandleEmptyDrawPileInAttackChain_WhenAttackAmountIsLarge tests reshuffling during attack chain
func TestShouldHandleEmptyDrawPileInAttackChain_WhenAttackAmountIsLarge(t *testing.T) {
	// Arrange
	player1 := player.New("player1")

	// Create a state with nearly empty draw pile, in an attack chain
	drawPile := deck.New()
	// Remove most cards from draw pile, leaving only 2
	for i := 0; i < 52; i++ {
		drawPile.Draw()
	}
	// Add specific cards using the correct method
	threeOfClubs := card.NewCard(card.Clubs, card.Three)
	fourOfHearts := card.NewCard(card.Hearts, card.Four)
	drawPile.AddToBottom(threeOfClubs)
	drawPile.AddToBottom(fourOfHearts)
	// Only 2 cards in draw pile

	discardPile := []card.Card{
		card.NewCard(card.Hearts, card.Five),
		card.NewCard(card.Diamonds, card.Six),
		card.NewCard(card.Clubs, card.Seven),
		card.NewCard(card.Spades, card.Seven), // Top card is a Seven (attack card)
	}

	topCard := discardPile[len(discardPile)-1]

	gameState := &state.State{
		Players:         []*player.Player{player1},
		ActivePlayers:   []string{player1.ID},
		CurrentPlayerId: player1.ID,
		Direction:       state.Clockwise,
		DrawPile:        drawPile,
		DiscardPile:     discardPile,
		TopCard:         topCard,
		InAttackChain:   true,
		AttackAmount:    6, // Player needs to draw 6 cards, but draw pile only has 2
		LastActiveSuit:  card.Spades,
	}

	initialHandSize := len(player1.Hand)

	// Act
	err := gameState.DrawCard("player1")

	// Assert
	if err != nil {
		t.Fatalf("Failed to draw cards during attack chain: %v", err)
	}

	// Check that player received all 6 cards from attack
	if len(player1.Hand) != initialHandSize+6 {
		t.Errorf("Expected player to draw 6 cards, hand size increased by %d", len(player1.Hand)-initialHandSize)
	}

	// Check that draw pile was reshuffled
	if gameState.DrawPile.IsEmpty() {
		t.Error("Expected draw pile to be refilled after reshuffling")
	}

	// Check that discard pile only contains top card
	if len(gameState.DiscardPile) != 1 {
		t.Errorf("Expected discard pile to only contain top card, got %d cards", len(gameState.DiscardPile))
	}
}
