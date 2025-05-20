package state

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/deck"
	"github.com/djoufson/check-games-engine/player"
)

func TestNewGame(t *testing.T) {
	playerIDs := []string{"player1", "player2", "player3"}

	// Test with default options
	state, err := New(playerIDs, nil)
	if err != nil {
		t.Fatalf("Failed to create new game: %v", err)
	}

	// Check basic game state
	if len(state.ActivePlayers) != 3 {
		t.Errorf("Expected 3 players, got %d", len(state.ActivePlayers))
	}

	if state.CurrentPlayerID() != "player1" {
		t.Errorf("Expected current player to be player1, got %s", state.CurrentPlayerID())
	}

	if state.IsGameOver() {
		t.Error("New game should not be over")
	}

	// Test with custom options
	opts := &GameOptions{
		InitialCards: 5,
		RandomSeed:   12345,
	}

	state, err = New(playerIDs, opts)
	if err != nil {
		t.Fatalf("Failed to create new game with options: %v", err)
	}

	// Check that player hands have expected size
	player1 := state.FindPlayerByID("player1")
	if player1 == nil {
		t.Fatalf("Failed to find player1")
	}

	if len(player1.Hand) != 5 {
		t.Errorf("Expected hand size to be 5, got %d", len(player1.Hand))
	}

	// Test with invalid player count
	_, err = New([]string{"player1"}, nil)
	if err == nil {
		t.Error("Expected error with only one player")
	}
}

func TestPlayCard(t *testing.T) {
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
	state := &State{
		Players:            []*player.Player{player1, player2},
		ActivePlayers:      []string{"player1", "player2"},
		CurrentPlayerIndex: 0,
		Direction:          Clockwise,
		DrawPile:           nil, // We won't be drawing in this test
		DiscardPile:        []card.Card{card.NewCard(card.Spades, card.Queen)},
		TopCard:            card.NewCard(card.Spades, card.Queen),
		InAttackChain:      false,
		AttackAmount:       0,
		LastActiveSuit:     card.Spades,
	}

	// Play a card that matches suit with the top card
	err := state.PlayCard("player1", card.NewCard(card.Spades, card.Ace))
	if err != nil {
		t.Fatalf("Failed to play matching card: %v", err)
	}

	// Check that the top card changed
	if state.TopCard.Suit != card.Spades || state.TopCard.Rank != card.Ace {
		t.Errorf("Expected top card to be Ace of Spades, got %v", state.TopCard)
	}

	// Check that it's player2's turn now (Ace skips, but with 2 players it goes back to player1)
	// But since Ace skips the next player, and there are only 2 players, it should be player1's turn again
	if state.CurrentPlayerID() != "player1" {
		t.Errorf("Expected current player to be player1 (after Ace skip), got %s", state.CurrentPlayerID())
	}

	// Check that player1's hand no longer contains the Ace
	if player1.HasCard(card.NewCard(card.Spades, card.Ace)) {
		t.Error("Expected player1's hand to no longer contain Ace of Spades")
	}
}

func TestDrawCard(t *testing.T) {
	// Set up a game state for testing DrawCard
	player1 := player.New("player1")
	player2 := player.New("player2")

	// Create a new deck with default cards
	drawPile := deck.New()

	// Get the top card for the initial discard pile
	topCard, _ := drawPile.Draw()

	state := &State{
		Players:            []*player.Player{player1, player2},
		ActivePlayers:      []string{"player1", "player2"},
		CurrentPlayerIndex: 0,
		Direction:          Clockwise,
		DrawPile:           drawPile,
		DiscardPile:        []card.Card{topCard},
		TopCard:            topCard,
		InAttackChain:      false,
		AttackAmount:       0,
		LastActiveSuit:     topCard.Suit,
	}

	// Get initial hand size
	initialHandSize := len(player1.Hand)

	// Draw a card
	err := state.DrawCard("player1")
	if err != nil {
		t.Fatalf("Failed to draw card: %v", err)
	}

	// Check that player1's hand has one more card
	if len(player1.Hand) != initialHandSize+1 {
		t.Errorf("Expected hand size to increase by 1, got %d", len(player1.Hand))
	}

	// Check that it's player2's turn now
	if state.CurrentPlayerID() != "player2" {
		t.Errorf("Expected current player to be player2, got %s", state.CurrentPlayerID())
	}
}

func TestAttackChain(t *testing.T) {
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

	state := &State{
		Players:            []*player.Player{player1, player2},
		ActivePlayers:      []string{"player1", "player2"},
		CurrentPlayerIndex: 0,
		Direction:          Clockwise,
		DrawPile:           drawPile,
		DiscardPile:        []card.Card{topCard},
		TopCard:            topCard,
		InAttackChain:      false,
		AttackAmount:       0,
		LastActiveSuit:     card.Hearts,
	}

	// Play a wild card to start an attack chain
	err := state.PlayCard("player1", card.NewCard(card.Hearts, card.Seven))
	if err != nil {
		t.Fatalf("Failed to play Seven: %v", err)
	}

	// Check attack chain started
	if !state.InAttackChain {
		t.Error("Expected to be in attack chain")
	}

	if state.AttackAmount != 2 {
		t.Errorf("Expected attack amount to be 2, got %d", state.AttackAmount)
	}

	// Check that it's player2's turn now
	if state.CurrentPlayerID() != "player2" {
		t.Errorf("Expected current player to be player2, got %s", state.CurrentPlayerID())
	}

	// Player2 responds with another wild card
	err = state.PlayCard("player2", card.NewCard(card.Spades, card.Seven))
	if err != nil {
		t.Fatalf("Failed to play Seven to continue attack chain: %v", err)
	}

	// Check attack amount increased
	if state.AttackAmount != 4 {
		t.Errorf("Expected attack amount to be 4, got %d", state.AttackAmount)
	}

	// Check that it's player1's turn now
	if state.CurrentPlayerID() != "player1" {
		t.Errorf("Expected current player to be player1, got %s", state.CurrentPlayerID())
	}

	// Player1 plays Joker to continue attack chain
	err = state.PlayCard("player1", card.NewRedJoker())
	if err != nil {
		t.Fatalf("Failed to play Joker: %v", err)
	}

	// Check attack amount increased
	if state.AttackAmount != 8 {
		t.Errorf("Expected attack amount to be 8, got %d", state.AttackAmount)
	}

	// Get initial hand size before player2 draws
	initialHandSize := len(player2.Hand)

	// Player2 has no more wild cards, so must draw
	err = state.DrawCard("player2")
	if err != nil {
		t.Fatalf("Failed to draw cards: %v", err)
	}

	// Check player2 drew 8 cards
	if len(player2.Hand) != initialHandSize+8 {
		t.Errorf("Expected hand size to increase by 8, got %d", len(player2.Hand))
	}

	// Check attack chain ended
	if state.InAttackChain {
		t.Error("Expected attack chain to end")
	}

	if state.AttackAmount != 0 {
		t.Errorf("Expected attack amount to reset to 0, got %d", state.AttackAmount)
	}
}

func TestIsGameOver(t *testing.T) {
	// Create players
	player1 := player.New("player1")
	player2 := player.New("player2")

	// Player1 has no cards (has won)
	player2.AddCardsToHand([]card.Card{
		card.NewCard(card.Spades, card.Ace),
	})

	state := &State{
		Players:            []*player.Player{player1, player2},
		ActivePlayers:      []string{"player2"}, // Only player2 is active
		CurrentPlayerIndex: 0,
		Direction:          Clockwise,
		DrawPile:           nil,
		DiscardPile:        []card.Card{card.NewCard(card.Hearts, card.King)},
		TopCard:            card.NewCard(card.Hearts, card.King),
		InAttackChain:      false,
		AttackAmount:       0,
		LastActiveSuit:     card.Hearts,
	}

	// Check game over
	if !state.IsGameOver() {
		t.Error("Expected game to be over")
	}

	// Check winner
	winners := state.GetWinner()
	if len(winners) != 1 || winners[0] != "player1" {
		t.Errorf("Expected player1 to be winner, got %v", winners)
	}

	// Check loser
	loser := state.GetLoser()
	if loser != "player2" {
		t.Errorf("Expected player2 to be loser, got %s", loser)
	}
}
