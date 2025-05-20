// Package game provides a high-level API for the check-game engine.
package game

import (
	"errors"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/player"
	"github.com/djoufson/check-games-engine/state"
)

// Game represents a check-game
type Game struct {
	state *state.State
}

// Options defines configurable options for a new game
type Options struct {
	InitialCards int   // Number of cards dealt to each player at start
	RandomSeed   int64 // Seed for RNG (useful for deterministic tests)
}

// State returns a snapshot of the current game state for serialization
func (g *Game) State() *state.State {
	return g.state.Clone()
}

// New creates a new game with the given player IDs and options
func New(playerIDs []string, options *Options) (*Game, error) {
	if len(playerIDs) < 2 {
		return nil, errors.New("at least 2 players are required")
	}

	var stateOpts *state.GameOptions
	if options != nil {
		stateOpts = &state.GameOptions{
			InitialCards: options.InitialCards,
			RandomSeed:   options.RandomSeed,
		}
	}

	gameState, err := state.New(playerIDs, stateOpts)
	if err != nil {
		return nil, err
	}

	return &Game{state: gameState}, nil
}

// FromState creates a game from an existing state
func FromState(s *state.State) *Game {
	return &Game{state: s}
}

// FromJSON creates a game from a JSON-serialized state
func FromJSON(data []byte) (*Game, error) {
	s, err := state.FromJSON(data)
	if err != nil {
		return nil, err
	}
	return &Game{state: s}, nil
}

// ToJSON serializes the game state to JSON
func (g *Game) ToJSON() ([]byte, error) {
	return g.state.ToJSON()
}

// CurrentPlayerID returns the ID of the player whose turn it is
func (g *Game) CurrentPlayerID() string {
	return g.state.CurrentPlayerID()
}

// IsPlayerTurn checks if it's the given player's turn
func (g *Game) IsPlayerTurn(playerID string) bool {
	return g.state.CurrentPlayerID() == playerID
}

// PlayCard plays a card from the current player's hand
func (g *Game) PlayCard(playerID string, c card.Card) error {
	return g.state.PlayCard(playerID, c)
}

// DrawCard causes the current player to draw a card
func (g *Game) DrawCard(playerID string) error {
	return g.state.DrawCard(playerID)
}

// ChangeSuit changes the active suit (used after playing a Jack)
func (g *Game) ChangeSuit(playerID string, newSuit card.Suit) error {
	return g.state.ChangeSuit(playerID, newSuit)
}

// GetPlayerHand returns the cards in the specified player's hand
func (g *Game) GetPlayerHand(playerID string) ([]card.Card, error) {
	player := g.state.FindPlayerByID(playerID)
	if player == nil {
		return nil, errors.New("player not found")
	}

	// Return a copy of the hand to prevent modification
	hand := make([]card.Card, len(player.Hand))
	copy(hand, player.Hand)

	return hand, nil
}

// GetPlayableCards returns the cards in the player's hand that can be played
func (g *Game) GetPlayableCards(playerID string) ([]card.Card, error) {
	if !g.IsPlayerTurn(playerID) {
		return nil, errors.New("not player's turn")
	}

	player := g.state.FindPlayerByID(playerID)
	if player == nil {
		return nil, errors.New("player not found")
	}

	return player.GetPlayableCards(g.state.TopCard, g.state.InAttackChain), nil
}

// GetTopCard returns the current top card
func (g *Game) GetTopCard() card.Card {
	return g.state.TopCard
}

// GetLastActiveSuit returns the last active suit (important for Jack effects)
func (g *Game) GetLastActiveSuit() card.Suit {
	return g.state.LastActiveSuit
}

// IsGameOver checks if the game is over
func (g *Game) IsGameOver() bool {
	return g.state.IsGameOver()
}

// GetWinners returns the IDs of players who have won
func (g *Game) GetWinners() []string {
	return g.state.GetWinner()
}

// GetLoser returns the ID of the player who lost
func (g *Game) GetLoser() string {
	return g.state.GetLoser()
}

// IsPlayerActive checks if a player is still in the game
func (g *Game) IsPlayerActive(playerID string) bool {
	return g.state.IsPlayerActive(playerID)
}

// IsInAttackChain checks if the game is in an attack chain
func (g *Game) IsInAttackChain() bool {
	return g.state.InAttackChain
}

// GetAttackAmount returns the current attack amount
func (g *Game) GetAttackAmount() int {
	return g.state.AttackAmount
}

// GetPlayerCount returns the number of players still in the game
func (g *Game) GetPlayerCount() int {
	return len(g.state.ActivePlayers)
}

// ValidateMove checks if a move is valid without modifying the game state
func (g *Game) ValidateMove(playerID string, c card.Card) (bool, error) {
	// Check if it's the player's turn
	if !g.IsPlayerTurn(playerID) {
		return false, errors.New("not player's turn")
	}

	// Find the player
	p := g.state.FindPlayerByID(playerID)
	if p == nil {
		return false, errors.New("player not found")
	}

	// Check if the player has the card
	if !p.HasCard(c) {
		return false, errors.New("card not in hand")
	}

	// Check if the play is valid according to game rules
	if !player.CanPlayCardOn(c, g.state.TopCard, g.state.InAttackChain) {
		return false, errors.New("invalid move")
	}

	return true, nil
}
