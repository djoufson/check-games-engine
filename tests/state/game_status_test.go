package state_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/player"
	"github.com/djoufson/check-games-engine/state"
)

func TestGameStatus(t *testing.T) {
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

	// Check game over
	if !gameState.IsGameOver() {
		t.Error("Expected game to be over")
	}

	// Check winner
	winners := gameState.GetWinner()
	if len(winners) != 1 || winners[0] != "player1" {
		t.Errorf("Expected player1 to be winner, got %v", winners)
	}

	// Check loser
	loser := gameState.GetLoser()
	if loser != "player2" {
		t.Errorf("Expected player2 to be loser, got %s", loser)
	}
}
