package state_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/deck"
	"github.com/djoufson/check-games-engine/player"
	"github.com/djoufson/check-games-engine/state"
)

// setupAttackChainTest creates a game state for testing attack chains
func setupAttackChainTest() (*state.State, *player.Player, *player.Player) {
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
		Players:         []*player.Player{player1, player2},
		ActivePlayers:   []string{player1.ID, player2.ID},
		CurrentPlayerId: player1.ID,
		Direction:       state.Clockwise,
		DrawPile:        drawPile,
		DiscardPile:     []card.Card{topCard},
		TopCard:         topCard,
		InAttackChain:   false,
		AttackAmount:    0,
		LastActiveSuit:  card.Hearts,
	}

	return gameState, player1, player2
}

// TestShouldStartAttackChain_WhenPlayingWildCard tests starting an attack chain
func TestShouldStartAttackChain_WhenPlayingWildCard(t *testing.T) {
	// Arrange
	gameState, _, _ := setupAttackChainTest()

	// Act
	err := gameState.PlayCard("player1", card.NewCard(card.Hearts, card.Seven))

	// Assert
	if err != nil {
		t.Fatalf("Failed to play Seven: %v", err)
	}

	if !gameState.InAttackChain {
		t.Error("Expected to be in attack chain")
	}

	if gameState.AttackAmount != 2 {
		t.Errorf("Expected attack amount to be 2, got %d", gameState.AttackAmount)
	}
}

// TestShouldSwitchTurns_WhenPlayingAttackCard tests turn switching after playing attack card
func TestShouldSwitchTurns_WhenPlayingAttackCard(t *testing.T) {
	// Arrange
	gameState, _, _ := setupAttackChainTest()

	// Act
	err := gameState.PlayCard("player1", card.NewCard(card.Hearts, card.Seven))

	// Assert
	if err != nil {
		t.Fatalf("Failed to play Seven: %v", err)
	}

	if gameState.CurrentPlayerID() != "player2" {
		t.Errorf("Expected current player to be player2, got %s", gameState.CurrentPlayerID())
	}
}

// TestShouldIncreaseAttackAmount_WhenRespondingWithWildCard tests attack amount increases
func TestShouldIncreaseAttackAmount_WhenRespondingWithWildCard(t *testing.T) {
	// Arrange
	gameState, _, _ := setupAttackChainTest()
	gameState.PlayCard("player1", card.NewCard(card.Hearts, card.Seven))

	// Act
	err := gameState.PlayCard("player2", card.NewCard(card.Spades, card.Seven))

	// Assert
	if err != nil {
		t.Fatalf("Failed to play Seven to continue attack chain: %v", err)
	}

	if gameState.AttackAmount != 4 {
		t.Errorf("Expected attack amount to be 4, got %d", gameState.AttackAmount)
	}
}

// TestShouldIncreaseAttackAmount_WhenPlayingJokerInAttackChain tests joker increases attack
func TestShouldIncreaseAttackAmount_WhenPlayingJokerInAttackChain(t *testing.T) {
	// Arrange
	gameState, _, _ := setupAttackChainTest()
	gameState.PlayCard("player1", card.NewCard(card.Hearts, card.Seven))
	gameState.PlayCard("player2", card.NewCard(card.Spades, card.Seven))

	// Act
	err := gameState.PlayCard("player1", card.NewRedJoker())

	// Assert
	if err != nil {
		t.Fatalf("Failed to play Joker: %v", err)
	}

	if gameState.AttackAmount != 8 {
		t.Errorf("Expected attack amount to be 8, got %d", gameState.AttackAmount)
	}
}

// TestShouldDrawPenaltyCards_WhenCannotRespondToAttack tests drawing penalty cards
func TestShouldDrawPenaltyCards_WhenCannotRespondToAttack(t *testing.T) {
	// Arrange
	gameState, _, player2 := setupAttackChainTest()
	gameState.PlayCard("player1", card.NewCard(card.Hearts, card.Seven))
	gameState.PlayCard("player2", card.NewCard(card.Spades, card.Seven))
	gameState.PlayCard("player1", card.NewRedJoker())

	initialHandSize := len(player2.Hand)

	// Act
	err := gameState.DrawCard("player2")

	// Assert
	if err != nil {
		t.Fatalf("Failed to draw cards: %v", err)
	}

	if len(player2.Hand) != initialHandSize+8 {
		t.Errorf("Expected hand size to increase by 8, got %d", len(player2.Hand))
	}
}

// TestShouldEndAttackChain_WhenDrawingPenaltyCards tests that attack chain ends after drawing
func TestShouldEndAttackChain_WhenDrawingPenaltyCards(t *testing.T) {
	// Arrange
	gameState, _, _ := setupAttackChainTest()
	gameState.PlayCard("player1", card.NewCard(card.Hearts, card.Seven))
	gameState.PlayCard("player2", card.NewCard(card.Spades, card.Seven))
	gameState.PlayCard("player1", card.NewRedJoker())

	// Act
	gameState.DrawCard("player2")

	// Assert
	if gameState.InAttackChain {
		t.Error("Expected attack chain to end")
	}
}
