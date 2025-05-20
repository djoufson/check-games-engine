package state_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/deck"
	"github.com/djoufson/check-games-engine/player"
	"github.com/djoufson/check-games-engine/state"
)

// setupCardPlayTest creates a test scenario for playing cards
func setupCardPlayTest() (*state.State, *player.Player, *player.Player) {
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

	return gameState, player1, player2
}

// setupCardDrawingTest creates a test scenario for drawing cards
func setupCardDrawingTest() (*state.State, *player.Player, *player.Player) {
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

	return gameState, player1, player2
}

// TestShouldChangeTopCard_WhenPlayingCard tests that the top card changes after playing
func TestShouldChangeTopCard_WhenPlayingCard(t *testing.T) {
	// Arrange
	gameState, _, _ := setupCardPlayTest()

	// Act
	err := gameState.PlayCard("player1", card.NewCard(card.Spades, card.Ace))

	// Assert
	if err != nil {
		t.Fatalf("Failed to play matching card: %v", err)
	}

	if gameState.TopCard.Suit != card.Spades || gameState.TopCard.Rank != card.Ace {
		t.Errorf("Expected top card to be Ace of Spades, got %v", gameState.TopCard)
	}
}

// TestShouldStayWithSamePlayer_WhenPlayingSkipCardWithTwoPlayers tests skip card with two players
func TestShouldStayWithSamePlayer_WhenPlayingSkipCardWithTwoPlayers(t *testing.T) {
	// Arrange
	gameState, _, _ := setupCardPlayTest()

	// Act
	err := gameState.PlayCard("player1", card.NewCard(card.Spades, card.Ace))

	// Assert
	if err != nil {
		t.Fatalf("Failed to play matching card: %v", err)
	}

	// With 2 players, if player1 plays a skip card, it should still be player1's turn
	if gameState.CurrentPlayerID() != "player1" {
		t.Errorf("Expected current player to be player1 (after Ace skip), got %s", gameState.CurrentPlayerID())
	}
}

// TestShouldRemoveCardFromHand_WhenPlayingCard tests that the card is removed from hand
func TestShouldRemoveCardFromHand_WhenPlayingCard(t *testing.T) {
	// Arrange
	gameState, player1, _ := setupCardPlayTest()

	// Act
	err := gameState.PlayCard("player1", card.NewCard(card.Spades, card.Ace))

	// Assert
	if err != nil {
		t.Fatalf("Failed to play matching card: %v", err)
	}

	if player1.HasCard(card.NewCard(card.Spades, card.Ace)) {
		t.Error("Expected player1's hand to no longer contain Ace of Spades")
	}
}

// TestShouldIncreaseHandSize_WhenDrawingCard tests that drawing increases hand size
func TestShouldIncreaseHandSize_WhenDrawingCard(t *testing.T) {
	// Arrange
	gameState, player1, _ := setupCardDrawingTest()
	initialHandSize := len(player1.Hand)

	// Act
	err := gameState.DrawCard("player1")

	// Assert
	if err != nil {
		t.Fatalf("Failed to draw card: %v", err)
	}

	if len(player1.Hand) != initialHandSize+1 {
		t.Errorf("Expected hand size to increase by 1, got %d", len(player1.Hand))
	}
}

// TestShouldChangeActivePlayer_WhenDrawingCard tests that drawing changes the active player
func TestShouldChangeActivePlayer_WhenDrawingCard(t *testing.T) {
	// Arrange
	gameState, _, _ := setupCardDrawingTest()

	// Act
	err := gameState.DrawCard("player1")

	// Assert
	if err != nil {
		t.Fatalf("Failed to draw card: %v", err)
	}

	if gameState.CurrentPlayerID() != "player2" {
		t.Errorf("Expected current player to be player2, got %s", gameState.CurrentPlayerID())
	}
}
