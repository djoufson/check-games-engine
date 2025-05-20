package game

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/deck"
	"github.com/djoufson/check-games-engine/player"
	"github.com/djoufson/check-games-engine/state"
)

func TestNewGame(t *testing.T) {
	playerIDs := []string{"player1", "player2", "player3"}

	// Test with default options
	game, err := New(playerIDs, nil)
	if err != nil {
		t.Fatalf("Failed to create new game: %v", err)
	}

	// Check basic game state
	if game.GetPlayerCount() != 3 {
		t.Errorf("Expected 3 players, got %d", game.GetPlayerCount())
	}

	if game.CurrentPlayerID() != "player1" {
		t.Errorf("Expected current player to be player1, got %s", game.CurrentPlayerID())
	}

	if game.IsGameOver() {
		t.Error("New game should not be over")
	}

	// Test with custom options
	opts := &Options{
		InitialCards: 5,
		RandomSeed:   12345,
	}

	game, err = New(playerIDs, opts)
	if err != nil {
		t.Fatalf("Failed to create new game with options: %v", err)
	}

	// Check that player hands have expected size
	hand, err := game.GetPlayerHand("player1")
	if err != nil {
		t.Fatalf("Failed to get player hand: %v", err)
	}

	if len(hand) != 5 {
		t.Errorf("Expected hand size to be 5, got %d", len(hand))
	}

	// Test with invalid player count
	_, err = New([]string{"player1"}, nil)
	if err == nil {
		t.Error("Expected error with only one player")
	}
}

func TestGameplayFlow(t *testing.T) {
	// Create a game with 2 players and a fixed seed for deterministic tests
	playerIDs := []string{"player1", "player2"}
	opts := &Options{
		InitialCards: 7,
		RandomSeed:   12345,
	}

	game, err := New(playerIDs, opts)
	if err != nil {
		t.Fatalf("Failed to create new game: %v", err)
	}

	// Get player1's hand
	player1Hand, err := game.GetPlayerHand("player1")
	if err != nil {
		t.Fatalf("Failed to get player1's hand: %v", err)
	}

	// Get the top card
	topCard := game.GetTopCard()

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
		err = game.DrawCard("player1")
		if err != nil {
			t.Fatalf("Failed to draw card: %v", err)
		}

		// Should be player2's turn now
		if game.CurrentPlayerID() != "player2" {
			t.Errorf("Expected current player to be player2, got %s", game.CurrentPlayerID())
		}

		// Get player2's hand
		player2Hand, err := game.GetPlayerHand("player2")
		if err != nil {
			t.Fatalf("Failed to get player2's hand: %v", err)
		}

		// Try to find a playable card in player2's hand
		foundPlayable = false
		for _, c := range player2Hand {
			if isValid, _ := game.ValidateMove("player2", c); isValid {
				cardToPlay = c
				foundPlayable = true
				break
			}
		}

		if foundPlayable {
			err = game.PlayCard("player2", cardToPlay)
			if err != nil {
				t.Fatalf("Failed to play card: %v", err)
			}

			// Should be player1's turn again
			if game.CurrentPlayerID() != "player1" {
				t.Errorf("Expected current player to be player1, got %s", game.CurrentPlayerID())
			}
		}
	} else {
		// Play the found card
		err = game.PlayCard("player1", cardToPlay)
		if err != nil {
			t.Fatalf("Failed to play card: %v", err)
		}

		// Check that the top card changed
		newTopCard := game.GetTopCard()
		if newTopCard.Suit != cardToPlay.Suit || newTopCard.Rank != cardToPlay.Rank {
			t.Errorf("Expected top card to match played card, got %v", newTopCard)
		}

		// Should be player2's turn now
		if game.CurrentPlayerID() != "player2" {
			t.Errorf("Expected current player to be player2, got %s", game.CurrentPlayerID())
		}
	}
}

func TestJSONSerialization(t *testing.T) {
	// Create a game
	playerIDs := []string{"player1", "player2"}
	opts := &Options{
		InitialCards: 7,
		RandomSeed:   12345,
	}

	game, _ := New(playerIDs, opts)

	// Serialize to JSON
	data, err := game.ToJSON()
	if err != nil {
		t.Fatalf("Failed to serialize game: %v", err)
	}

	// Create a new game from the JSON
	game2, err := FromJSON(data)
	if err != nil {
		t.Fatalf("Failed to deserialize game: %v", err)
	}

	// Check that the games match
	if game.CurrentPlayerID() != game2.CurrentPlayerID() {
		t.Errorf("Current player mismatch: %s vs %s",
			game.CurrentPlayerID(), game2.CurrentPlayerID())
	}

	if game.GetPlayerCount() != game2.GetPlayerCount() {
		t.Errorf("Player count mismatch: %d vs %d",
			game.GetPlayerCount(), game2.GetPlayerCount())
	}

	// Check that the top cards match
	topCard1 := game.GetTopCard()
	topCard2 := game2.GetTopCard()

	if topCard1.Suit != topCard2.Suit || topCard1.Rank != topCard2.Rank {
		t.Errorf("Top card mismatch: %v vs %v", topCard1, topCard2)
	}
}

func TestSpecialCardEffects(t *testing.T) {
	// Set up players with specific cards for testing
	jackCard := card.NewCard(card.Spades, card.Jack)   // Suit change - same suit as top card
	aceCard := card.NewCard(card.Spades, card.Ace)     // Skip
	sevenCard := card.NewCard(card.Hearts, card.Seven) // Wild (+2)
	jokerCard := card.NewRedJoker()                    // Wild (+4)

	player1 := player.New("player1")
	player1.AddCardsToHand([]card.Card{
		aceCard,
		sevenCard,
		jackCard,
		jokerCard,
	})

	player2 := player.New("player2")
	player2.AddCardsToHand([]card.Card{
		card.NewCard(card.Hearts, card.Five),
		card.NewCard(card.Spades, card.Seven), // Wild (+2)
		card.NewCard(card.Clubs, card.Jack),   // Suit change
		card.NewBlackJoker(),                  // Wild (+4)
	})

	// Create a state with these players and specific top card
	stateObj := &state.State{
		Players:            []*player.Player{player1, player2},
		ActivePlayers:      []string{player1.ID, player2.ID},
		CurrentPlayerIndex: 0,
		Direction:          state.Clockwise,
		DrawPile:           deck.New(), // Standard deck
		DiscardPile:        []card.Card{card.NewCard(card.Spades, card.Queen)},
		TopCard:            card.NewCard(card.Spades, card.Queen),
		InAttackChain:      false,
		AttackAmount:       0,
		LastActiveSuit:     card.Spades,
	}

	// Initialize a game from this state
	game := FromState(stateObj)

	// Test Ace (Skip)
	err := game.PlayCard(player1.ID, aceCard)
	if err != nil {
		t.Fatalf("Failed to play Ace: %v", err)
	}

	// With 2 players, it should stay player1's turn
	if game.CurrentPlayerID() != "player1" {
		t.Errorf("Expected current player to be player1 after Ace, got %s", game.CurrentPlayerID())
	}

	// Test Jack (Suit change)
	// First play Jack
	err = game.PlayCard(player1.ID, jackCard)
	if err != nil {
		t.Fatalf("Failed to play Jack: %v", err)
	}

	// Then change suit
	err = game.ChangeSuit(player1.ID, card.Hearts)
	if err != nil {
		t.Fatalf("Failed to change suit: %v", err)
	}

	// LastActiveSuit should be Hearts now
	if game.GetLastActiveSuit() != card.Hearts {
		t.Errorf("Expected last active suit to be Hearts, got %v", game.GetLastActiveSuit())
	}

	// Test wild card (Seven) - should be player2's turn now
	player2SevenCard := card.NewCard(card.Spades, card.Seven)
	err = game.PlayCard(player2.ID, player2SevenCard)
	if err != nil {
		t.Fatalf("Failed to play Seven: %v", err)
	}

	// Should be in attack chain
	if !game.IsInAttackChain() {
		t.Error("Expected to be in attack chain")
	}

	if game.GetAttackAmount() != 2 {
		t.Errorf("Expected attack amount to be 2, got %d", game.GetAttackAmount())
	}

	// Test wild card (Joker) in attack chain
	err = game.PlayCard(player1.ID, jokerCard)
	if err != nil {
		t.Fatalf("Failed to play Joker: %v", err)
	}

	// Attack amount should be increased
	if game.GetAttackAmount() != 6 {
		t.Errorf("Expected attack amount to be 6, got %d", game.GetAttackAmount())
	}
}
