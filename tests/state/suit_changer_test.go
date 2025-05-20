package state_test

import (
	"testing"

	"fmt"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/deck"
	"github.com/djoufson/check-games-engine/player"
	"github.com/djoufson/check-games-engine/state"
)

// setupSuitChangerTest creates a game state for testing suit changing with Jack
func setupSuitChangerTest() (*state.State, *player.Player, *player.Player) {
	// Create players
	player1 := player.New("player1")
	player1.AddCardsToHand([]card.Card{
		card.NewCard(card.Clubs, card.Jack), // Jack for changing suit
		card.NewCard(card.Spades, card.King),
	})

	player2 := player.New("player2")
	player2.AddCardsToHand([]card.Card{
		card.NewCard(card.Hearts, card.Five),
		card.NewCard(card.Diamonds, card.Six),
	})

	// Create a state with a known top card
	gameState := &state.State{
		Players:         []*player.Player{player1, player2},
		ActivePlayers:   []string{player1.ID, player2.ID},
		CurrentPlayerId: player1.ID,
		Direction:       state.Clockwise,
		DrawPile:        deck.New(),
		DiscardPile:     []card.Card{card.NewCard(card.Diamonds, card.Queen)},
		TopCard:         card.NewCard(card.Diamonds, card.Queen),
		InAttackChain:   false,
		AttackAmount:    0,
		LastActiveSuit:  card.Diamonds,
	}

	return gameState, player1, player2
}

// TestShouldChangeSuit_WhenPlayingJack tests that playing a Jack changes the active suit
func TestShouldChangeSuit_WhenPlayingJack(t *testing.T) {
	// Arrange
	gameState, _, _ := setupSuitChangerTest()
	jackCard := card.NewCard(card.Clubs, card.Jack)

	// Act
	// Play a Jack and change the suit to Spades
	err := gameState.PlayCard("player1", jackCard)
	if err != nil {
		t.Fatalf("Failed to play Jack: %v", err)
	}

	err = gameState.ChangeSuit("player1", card.Spades)

	// Assert
	if err != nil {
		t.Fatalf("Failed to change suit: %v", err)
	}

	if gameState.LastActiveSuit != card.Spades {
		t.Errorf("Expected active suit to be Spades, got %v", gameState.LastActiveSuit)
	}

	if gameState.TopCard.Suit != card.Clubs || gameState.TopCard.Rank != card.Jack {
		t.Errorf("Expected top card to be Jack of Clubs, got %v", gameState.TopCard)
	}
}

// TestShouldRequireValidSuit_WhenChangingSuit tests suit validation when changing suit
func TestShouldRequireValidSuit_WhenChangingSuit(t *testing.T) {
	// Arrange
	gameState, _, _ := setupSuitChangerTest()
	jackCard := card.NewCard(card.Clubs, card.Jack)

	// Act
	err := gameState.PlayCard("player1", jackCard)
	if err != nil {
		t.Fatalf("Failed to play Jack: %v", err)
	}

	// Try to set an invalid suit (outside the valid enum values)
	invalidSuit := card.Suit(fmt.Sprint(99))
	err = gameState.ChangeSuit("player1", invalidSuit)

	// Assert
	if err == nil {
		t.Error("Expected error when changing to invalid suit, got nil")
	}
}

// TestShouldEnforceSuitChange_WhenNextPlayerPlays tests that next player must follow the new suit
func TestShouldEnforceSuitChange_WhenNextPlayerPlays(t *testing.T) {
	// Arrange
	gameState, _, _ := setupSuitChangerTest()
	jackCard := card.NewCard(card.Clubs, card.Jack)

	// Player1 plays Jack and changes suit to Spades
	err := gameState.PlayCard("player1", jackCard)
	if err != nil {
		t.Fatalf("Failed to play Jack: %v", err)
	}

	err = gameState.ChangeSuit("player1", card.Spades)
	if err != nil {
		t.Fatalf("Failed to change suit: %v", err)
	}

	// Act
	// Player2 tries to play a Hearts card, which should fail
	heartsCard := card.NewCard(card.Hearts, card.Five)
	err = gameState.PlayCard("player2", heartsCard)

	// Assert
	if err == nil {
		t.Error("Expected error when playing card that doesn't match the changed suit")
	}

	// Act again
	// Player2 draws a card instead (should work)
	err = gameState.DrawCard("player2")

	// Assert
	if err != nil {
		t.Errorf("Failed to draw card after suit change: %v", err)
	}
}
