package state_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/player"
	"github.com/djoufson/check-games-engine/state"
)

// setupGameOverTest creates a game state scenario where the game is over
func setupGameOverTest() (*state.State, *player.Player, *player.Player) {
	// Create players
	player1 := player.New("player1")
	player2 := player.New("player2")

	// Player1 has no cards (has won)
	player2.AddCardsToHand([]card.Card{
		card.NewCard(card.Spades, card.Ace),
	})

	gameState := &state.State{
		Players:            []*player.Player{player1, player2},
		ActivePlayers:      []string{"player2"}, // Only player2 is active
		CurrentPlayerIndex: 0,
		Direction:          state.Clockwise,
		DrawPile:           nil,
		DiscardPile:        []card.Card{card.NewCard(card.Hearts, card.King)},
		TopCard:            card.NewCard(card.Hearts, card.King),
		InAttackChain:      false,
		AttackAmount:       0,
		LastActiveSuit:     card.Hearts,
	}

	return gameState, player1, player2
}

// TestShouldReportGameOver_WhenPlayerHasNoCards tests game over detection
func TestShouldReportGameOver_WhenPlayerHasNoCards(t *testing.T) {
	// Arrange
	gameState, _, _ := setupGameOverTest()

	// Act & Assert
	if !gameState.IsGameOver() {
		t.Error("Expected game to be over")
	}
}

// TestShouldIdentifyWinner_WhenPlayerHasNoCards tests winner detection
func TestShouldIdentifyWinner_WhenPlayerHasNoCards(t *testing.T) {
	// Arrange
	gameState, _, _ := setupGameOverTest()

	// Act
	winners := gameState.GetWinner()

	// Assert
	if len(winners) != 1 || winners[0] != "player1" {
		t.Errorf("Expected player1 to be winner, got %v", winners)
	}
}

// TestShouldIdentifyLoser_WhenGameIsOver tests loser detection
func TestShouldIdentifyLoser_WhenGameIsOver(t *testing.T) {
	// Arrange
	gameState, _, _ := setupGameOverTest()

	// Act
	loser := gameState.GetLoser()

	// Assert
	if loser != "player2" {
		t.Errorf("Expected player2 to be loser, got %s", loser)
	}
}
