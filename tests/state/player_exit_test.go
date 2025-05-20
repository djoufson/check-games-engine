package state_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/player"
	"github.com/djoufson/check-games-engine/state"
)

// setupLastCardTest creates a game state for testing playing the last card
func setupLastCardTest() (*state.State, *player.Player, *player.Player) {
	// Create players
	player1 := player.New("player1")
	player1.AddCardsToHand([]card.Card{
		card.NewCard(card.Hearts, card.King), // Player1's only card
	})

	player2 := player.New("player2")
	player2.AddCardsToHand([]card.Card{
		card.NewCard(card.Hearts, card.Queen),
		card.NewCard(card.Spades, card.Jack),
	})

	// Create a state with a hearts card on top
	gameState := &state.State{
		Players:         []*player.Player{player1, player2},
		ActivePlayers:   []string{player1.ID, player2.ID},
		CurrentPlayerId: player1.ID,
		Direction:       state.Clockwise,
		DrawPile:        nil, // Not important for this test
		DiscardPile:     []card.Card{card.NewCard(card.Hearts, card.Ten)},
		TopCard:         card.NewCard(card.Hearts, card.Ten),
		InAttackChain:   false,
		AttackAmount:    0,
		LastActiveSuit:  card.Hearts,
	}

	return gameState, player1, player2
}

// setupLastCardInAttackChainTest creates a test environment where a player can play their last card during an attack
func setupLastCardInAttackChainTest() (*state.State, *player.Player, *player.Player, *player.Player) {
	// Create three players
	player1 := player.New("player1")
	player1.AddCardsToHand([]card.Card{
		card.NewCard(card.Hearts, card.Seven), // Wild attack card
	})

	player2 := player.New("player2")
	player2.AddCardsToHand([]card.Card{
		card.NewCard(card.Spades, card.Seven), // Player2's only card - a response to the attack
	})

	player3 := player.New("player3")
	player3.AddCardsToHand([]card.Card{
		card.NewCard(card.Hearts, card.Queen),
		card.NewCard(card.Spades, card.Jack),
	})

	// Create a state with an attack chain started
	gameState := &state.State{
		Players:         []*player.Player{player1, player2, player3},
		ActivePlayers:   []string{player1.ID, player2.ID, player3.ID},
		CurrentPlayerId: player1.ID,
		Direction:       state.Clockwise,
		DrawPile:        nil, // Not important for this test
		DiscardPile:     []card.Card{card.NewCard(card.Hearts, card.Ten)},
		TopCard:         card.NewCard(card.Hearts, card.Ten),
		InAttackChain:   false,
		AttackAmount:    0,
		LastActiveSuit:  card.Hearts,
	}

	return gameState, player1, player2, player3
}

// TestShouldExitGame_WhenPlayingLastCard tests that a player exits when playing last card
func TestShouldExitGame_WhenPlayingLastCard(t *testing.T) {
	// Arrange
	gameState, _, _ := setupLastCardTest()

	// Act
	err := gameState.PlayCard("player1", card.NewCard(card.Hearts, card.King))

	// Assert
	if err != nil {
		t.Fatalf("Failed to play card: %v", err)
	}

	// Check that player1 is no longer active
	for _, activeID := range gameState.ActivePlayers {
		if activeID == "player1" {
			t.Error("Expected player1 to no longer be in active players list")
		}
	}

	// Verify player2 is now the current player
	if gameState.CurrentPlayerID() != "player2" {
		t.Errorf("Expected player2 to be current player, got %s", gameState.CurrentPlayerID())
	}
}

// TestShouldExitGame_WhenPlayingLastCardInAttackChain tests playing last card during attack chain
func TestShouldExitGame_WhenPlayingLastCardInAttackChain(t *testing.T) {
	// Arrange
	gameState, _, _, _ := setupLastCardInAttackChainTest()

	// Start attack chain
	err := gameState.PlayCard("player1", card.NewCard(card.Hearts, card.Seven))
	if err != nil {
		t.Fatalf("Failed to play Seven to start attack: %v", err)
	}

	// Verify attack chain started
	if !gameState.InAttackChain {
		t.Fatal("Expected to be in attack chain")
	}

	// Act
	// Player2 plays their last card (a Seven) to continue the attack
	err = gameState.PlayCard("player2", card.NewCard(card.Spades, card.Seven))

	// Assert
	if err != nil {
		t.Fatalf("Failed to play last card during attack chain: %v", err)
	}

	// Check that player2 is no longer active
	for _, activeID := range gameState.ActivePlayers {
		if activeID == "player2" {
			t.Error("Expected player2 to no longer be in active players list")
		}
	}

	// Verify attack chain continues to player3
	if !gameState.InAttackChain {
		t.Error("Expected attack chain to continue to next player")
	}

	// Verify attack amount increased
	if gameState.AttackAmount != 4 {
		t.Errorf("Expected attack amount to be 4, got %d", gameState.AttackAmount)
	}

	// Verify player3 is now the current player
	if gameState.CurrentPlayerID() != "player3" {
		t.Errorf("Expected player3 to be current player, got %s", gameState.CurrentPlayerID())
	}
}

// TestShouldDetermineWinner_WhenOnlyOnePlayerRemains tests winner determination when only one player remains
func TestShouldDetermineWinner_WhenOnlyOnePlayerRemains(t *testing.T) {
	// Arrange
	// player1 with one card and player2 with two cards
	gameState, _, _ := setupLastCardTest()

	// Act
	// Player1 plays their last card
	gameState.PlayCard("player1", card.NewCard(card.Hearts, card.King))

	// Assert
	// Game is over yet because player2 is alone in the game
	if !gameState.IsGameOver() {
		t.Error("Game should be over with one player left")
	}

	// player1 should be the winner (first to empty hand)
	winners := gameState.GetWinner()
	if len(winners) != 1 || winners[0] != "player1" {
		t.Errorf("Expected player1 to be winner, got %v", winners)
	}

	// player2 should be the loser (last with cards)
	loser := gameState.GetLoser()
	if loser != "player2" {
		t.Errorf("Expected player2 to be loser, got %s", loser)
	}
}
