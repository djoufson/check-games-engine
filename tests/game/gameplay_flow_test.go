package game_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/game"
)

func TestGameplayFlow(t *testing.T) {
	// Create a game with 2 players and a fixed seed for deterministic tests
	playerIDs := []string{"player1", "player2"}
	opts := &game.Options{
		InitialCards: 7,
		RandomSeed:   12345,
	}

	g, err := game.New(playerIDs, opts)
	if err != nil {
		t.Fatalf("Failed to create new game: %v", err)
	}

	// Get player1's hand
	player1Hand, err := g.GetPlayerHand("player1")
	if err != nil {
		t.Fatalf("Failed to get player1's hand: %v", err)
	}

	// Get the top card
	topCard := g.GetTopCard()

	// Find a playable card in player1's hand
	var cardToPlay card.Card
	var foundPlayable bool

	for _, c := range player1Hand {
		// Try to find a standard card (not special effect) just for simplicity in testing
		if !c.IsWildCard() && !c.IsSkip() && !c.IsSuitChanger() && !c.IsTransparent() {
			if c.Suit == topCard.Suit || c.Rank == topCard.Rank {
				cardToPlay = c
				foundPlayable = true
				break
			}
		}
	}

	// If no standard playable card found, player1 draws
	if !foundPlayable {
		err = g.DrawCard("player1")
		if err != nil {
			t.Fatalf("Failed to draw card: %v", err)
		}

		// Should be player2's turn now
		if g.CurrentPlayerID() != "player2" {
			t.Errorf("Expected current player to be player2, got %s", g.CurrentPlayerID())
		}

		// Get player2's hand
		player2Hand, err := g.GetPlayerHand("player2")
		if err != nil {
			t.Fatalf("Failed to get player2's hand: %v", err)
		}

		// Try to find a playable card in player2's hand
		foundPlayable = false
		for _, c := range player2Hand {
			if isValid, _ := g.ValidateMove("player2", c); isValid {
				cardToPlay = c
				foundPlayable = true
				break
			}
		}

		if foundPlayable {
			err = g.PlayCard("player2", cardToPlay)
			if err != nil {
				t.Fatalf("Failed to play card: %v", err)
			}

			// Should be player1's turn again
			if g.CurrentPlayerID() != "player1" {
				t.Errorf("Expected current player to be player1, got %s", g.CurrentPlayerID())
			}
		}
	} else {
		// Play the found card
		err = g.PlayCard("player1", cardToPlay)
		if err != nil {
			t.Fatalf("Failed to play card: %v", err)
		}

		// Check that the top card changed
		newTopCard := g.GetTopCard()
		if newTopCard.Suit != cardToPlay.Suit || newTopCard.Rank != cardToPlay.Rank {
			t.Errorf("Expected top card to match played card, got %v", newTopCard)
		}

		// Should be player2's turn now
		if g.CurrentPlayerID() != "player2" {
			t.Errorf("Expected current player to be player2, got %s", g.CurrentPlayerID())
		}
	}
}
