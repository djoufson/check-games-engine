package state_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/player"
	"github.com/djoufson/check-games-engine/state"
)

// setupAttackChainRestrictionsTest creates a test environment for attack chain restrictions
func setupAttackChainRestrictionsTest() (*state.State, *player.Player, *player.Player) {
	// Create players with specific cards
	player1 := player.New("player1")
	player1.AddCardsToHand([]card.Card{
		card.NewCard(card.Hearts, card.Seven), // Wild (+2)
	})

	player2 := player.New("player2")
	player2.AddCardsToHand([]card.Card{
		card.NewCard(card.Spades, card.Two),    // Transparent
		card.NewCard(card.Diamonds, card.Jack), // Suit changer
		card.NewCard(card.Hearts, card.King),   // Regular card
	})

	// Create a state with known top card and attack chain active
	topCard := card.NewCard(card.Hearts, card.Seven)

	gameState := &state.State{
		Players:         []*player.Player{player1, player2},
		ActivePlayers:   []string{player1.ID, player2.ID},
		CurrentPlayerId: player2.ID, // Player2's turn (needs to respond to attack)
		Direction:       state.Clockwise,
		DrawPile:        nil, // Not important for this test
		DiscardPile:     []card.Card{topCard},
		TopCard:         topCard,
		InAttackChain:   true,
		AttackAmount:    2, // 2 cards to draw (from the Seven)
		LastActiveSuit:  card.Hearts,
	}

	return gameState, player1, player2
}

// TestShouldRejectTransparentCard_WhenInAttackChain tests that transparent 2s can't be used in attack chains
func TestShouldRejectTransparentCard_WhenInAttackChain(t *testing.T) {
	// Arrange
	gameState, _, _ := setupAttackChainRestrictionsTest()

	// Act
	// Player2 tries to play a transparent 2, which should be rejected
	transparentCard := card.NewCard(card.Spades, card.Two)
	err := gameState.PlayCard("player2", transparentCard)

	// Assert
	if err == nil {
		t.Error("Expected error when playing transparent 2 in attack chain, got nil")
	}

	// Verify attack chain is still active
	if !gameState.InAttackChain {
		t.Error("Attack chain should still be active")
	}

	// Verify attack amount is unchanged
	if gameState.AttackAmount != 2 {
		t.Errorf("Expected attack amount to remain 2, got %d", gameState.AttackAmount)
	}
}

// TestShouldRejectSuitChanger_WhenInAttackChain tests that Jacks can't be used in attack chains
func TestShouldRejectSuitChanger_WhenInAttackChain(t *testing.T) {
	// Arrange
	gameState, _, _ := setupAttackChainRestrictionsTest()

	// Act
	// Player2 tries to play a Jack, which should be rejected
	jackCard := card.NewCard(card.Diamonds, card.Jack)
	err := gameState.PlayCard("player2", jackCard)

	// Assert
	if err == nil {
		t.Error("Expected error when playing Jack in attack chain, got nil")
	}

	// Verify attack chain is still active
	if !gameState.InAttackChain {
		t.Error("Attack chain should still be active")
	}

	// Verify attack amount is unchanged
	if gameState.AttackAmount != 2 {
		t.Errorf("Expected attack amount to remain 2, got %d", gameState.AttackAmount)
	}
}

// TestShouldAcceptWildCard_WhenInAttackChain tests that wild cards can be used in attack chains
func TestShouldAcceptWildCard_WhenInAttackChain(t *testing.T) {
	// Arrange
	gameState, _, player2 := setupAttackChainRestrictionsTest()

	// Add a wild card to player2's hand
	sevenCard := card.NewCard(card.Clubs, card.Seven)
	player2.AddCardsToHand([]card.Card{sevenCard})

	// Act
	// Player2 plays another Seven to continue the attack chain
	err := gameState.PlayCard("player2", sevenCard)

	// Assert
	if err != nil {
		t.Fatalf("Failed to play Seven in attack chain: %v", err)
	}

	// Verify attack chain is still active
	if !gameState.InAttackChain {
		t.Error("Attack chain should still be active")
	}

	// Verify attack amount increased
	if gameState.AttackAmount != 4 {
		t.Errorf("Expected attack amount to be 4, got %d", gameState.AttackAmount)
	}
}
