package game_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/game"
)

// setupGameplayTest creates a new game for testing gameplay flow
func setupGameplayTest(t *testing.T) *game.Game {
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

	return g
}

// findPlayableCard finds a standard playable card in a player's hand
func findPlayableCard(t *testing.T, g *game.Game, playerID string) (card.Card, bool) {
	playerHand, err := g.GetPlayerHand(playerID)
	if err != nil {
		t.Fatalf("Failed to get %s's hand: %v", playerID, err)
	}

	topCard := g.GetTopCard()

	// Find a playable card in the player's hand
	for _, c := range playerHand {
		// Try to find a standard card (not special effect) just for simplicity in testing
		if !c.IsWildCard() && !c.IsSkip() && !c.IsSuitChanger() && !c.IsTransparent() {
			if c.Suit == topCard.Suit || c.Rank == topCard.Rank {
				return c, true
			}
		}
	}

	// If no standard card found, check for any valid card
	for _, c := range playerHand {
		if isValid, _ := g.ValidateMove(playerID, c); isValid {
			return c, true
		}
	}

	return card.Card{}, false
}

// TestShouldSwitchTurns_WhenDrawingCard tests that drawing a card switches turns
func TestShouldSwitchTurns_WhenDrawingCard(t *testing.T) {
	// Arrange
	g := setupGameplayTest(t)

	// Act
	err := g.DrawCard("player1")

	// Assert
	if err != nil {
		t.Fatalf("Failed to draw card: %v", err)
	}

	if g.CurrentPlayerID() != "player2" {
		t.Errorf("Expected current player to be player2, got %s", g.CurrentPlayerID())
	}
}

// TestShouldUpdateTopCard_WhenPlayingCard tests that playing a card updates the top card
func TestShouldUpdateTopCard_WhenPlayingCard(t *testing.T) {
	// Arrange
	g := setupGameplayTest(t)
	cardToPlay, foundPlayable := findPlayableCard(t, g, "player1")
	if !foundPlayable {
		t.Skip("No playable card found in player1's hand")
	}

	// Act
	err := g.PlayCard("player1", cardToPlay)

	// Assert
	if err != nil {
		t.Fatalf("Failed to play card: %v", err)
	}

	newTopCard := g.GetTopCard()
	if newTopCard.Suit != cardToPlay.Suit || newTopCard.Rank != cardToPlay.Rank {
		t.Errorf("Expected top card to match played card, got %v", newTopCard)
	}
}

// TestShouldSwitchTurns_WhenPlayingCard tests that playing a card switches turns
func TestShouldSwitchTurns_WhenPlayingCard(t *testing.T) {
	// Arrange
	g := setupGameplayTest(t)
	cardToPlay, foundPlayable := findPlayableCard(t, g, "player1")
	if !foundPlayable {
		t.Skip("No playable card found in player1's hand")
	}

	// Act
	err := g.PlayCard("player1", cardToPlay)

	// Assert
	if err != nil {
		t.Fatalf("Failed to play card: %v", err)
	}

	// Should be player2's turn now (unless it was a special card)
	if !cardToPlay.IsSkip() && g.CurrentPlayerID() != "player2" {
		t.Errorf("Expected current player to be player2, got %s", g.CurrentPlayerID())
	}
}
