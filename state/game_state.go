// Package state defines the game state and state transitions for the check-game engine.
package state

import (
	"encoding/json"
	"errors"
	"math/rand"

	"github.com/check-games/engine/card"
	"github.com/check-games/engine/deck"
	"github.com/check-games/engine/player"
)

// Direction represents the direction of play
type Direction int

const (
	// Clockwise is the standard play direction
	Clockwise Direction = iota
	// CounterClockwise is the reverse play direction
	CounterClockwise
)

// State represents the current state of a game
type State struct {
	Players            []*player.Player `json:"players"`
	ActivePlayers      []string         `json:"active_players"` // IDs of players still in the game
	CurrentPlayerIndex int              `json:"current_player_index"`
	Direction          Direction        `json:"direction"`
	DrawPile           *deck.Deck       `json:"draw_pile"`
	DiscardPile        []card.Card      `json:"discard_pile"`
	TopCard            card.Card        `json:"top_card"`
	InAttackChain      bool             `json:"in_attack_chain"`
	AttackAmount       int              `json:"attack_amount"`
	LastActiveSuit     card.Suit        `json:"last_active_suit"` // For Jack's suit change effect
	LockedTurn         bool             `json:"blocked_turn"`     // If the turn is blocked until the suit is changed
}

// GameOptions defines configurable options for a new game
type GameOptions struct {
	InitialCards  int              // Number of cards dealt to each player at start
	RandomSeed    int64            // Seed for RNG (useful for deterministic tests)
	CustomPlayers []*player.Player // For testing or restarting a game
}

// DefaultOptions returns the default game options
func DefaultOptions() GameOptions {
	return GameOptions{
		InitialCards: 7,
		RandomSeed:   0, // Will use time if 0
	}
}

// New creates a new game state with the given player IDs and options
func New(playerIDs []string, options *GameOptions) (*State, error) {
	if len(playerIDs) < 2 {
		return nil, errors.New("at least 2 players are required")
	}

	// Use default options if none provided
	opts := DefaultOptions()
	if options != nil {
		opts = *options
	}

	// Create a new seed if none provided
	seed := opts.RandomSeed

	// Create a new RNG
	r := rand.New(rand.NewSource(seed))

	// Create and shuffle the deck
	drawPile := deck.New()
	drawPile.Shuffle(r)

	// Create players
	players := make([]*player.Player, len(playerIDs))
	activePlayerIDs := make([]string, len(playerIDs))

	for i, id := range playerIDs {
		players[i] = player.New(id)
		activePlayerIDs[i] = id
	}

	// Deal initial cards
	for i := 0; i < opts.InitialCards; i++ {
		for _, p := range players {
			if card, ok := drawPile.Draw(); ok {
				p.AddToHand(card)
			}
		}
	}

	// Draw the top card for the discard pile
	topCard, ok := drawPile.Draw()
	if !ok {
		return nil, errors.New("failed to draw initial card for discard pile")
	}

	// Ensure the initial discard card isn't a wild card
	if topCard.IsWildCard() {
		// Put the wild card back in the deck and shuffle again
		drawPile.AddToBottom(topCard)
		drawPile.Shuffle(r)

		// Draw again
		topCard, ok = drawPile.Draw()
		if !ok {
			return nil, errors.New("failed to draw initial card for discard pile")
		}
	}

	discardPile := []card.Card{topCard}

	// Create the initial game state
	state := &State{
		Players:            players,
		ActivePlayers:      activePlayerIDs,
		CurrentPlayerIndex: 0,
		Direction:          Clockwise,
		DrawPile:           drawPile,
		DiscardPile:        discardPile,
		TopCard:            topCard,
		InAttackChain:      false,
		AttackAmount:       0,
		LastActiveSuit:     topCard.Suit,
	}

	return state, nil
}

// Clone creates a deep copy of the game state
func (s *State) Clone() *State {
	// Clone players
	players := make([]*player.Player, len(s.Players))
	for i, p := range s.Players {
		clone := player.New(p.ID)
		clone.Hand = make([]card.Card, len(p.Hand))
		copy(clone.Hand, p.Hand)
		players[i] = clone
	}

	// Clone active players
	activePlayerIDs := make([]string, len(s.ActivePlayers))
	copy(activePlayerIDs, s.ActivePlayers)

	// Clone draw pile
	drawPile := &deck.Deck{
		Cards: make([]card.Card, len(s.DrawPile.Cards)),
	}
	copy(drawPile.Cards, s.DrawPile.Cards)

	// Clone discard pile
	discardPile := make([]card.Card, len(s.DiscardPile))
	copy(discardPile, s.DiscardPile)

	// Create the cloned state
	clone := &State{
		Players:            players,
		ActivePlayers:      activePlayerIDs,
		CurrentPlayerIndex: s.CurrentPlayerIndex,
		Direction:          s.Direction,
		DrawPile:           drawPile,
		DiscardPile:        discardPile,
		TopCard:            s.TopCard,
		InAttackChain:      s.InAttackChain,
		AttackAmount:       s.AttackAmount,
		LastActiveSuit:     s.LastActiveSuit,
	}

	return clone
}

// CurrentPlayerID returns the ID of the player whose turn it is
func (s *State) CurrentPlayerID() string {
	return s.ActivePlayers[s.CurrentPlayerIndex]
}

// CurrentPlayer returns the player whose turn it is
func (s *State) CurrentPlayer() *player.Player {
	currentID := s.CurrentPlayerID()
	for _, p := range s.Players {
		if p.ID == currentID {
			return p
		}
	}
	return nil
}

// NextPlayerIndex computes the index of the next player based on current direction
func (s *State) NextPlayerIndex() int {
	numActivePlayers := len(s.ActivePlayers)
	if numActivePlayers <= 1 {
		return 0 // Only one player left
	}

	nextIdx := s.CurrentPlayerIndex
	if s.Direction == Clockwise {
		nextIdx = (nextIdx + 1) % numActivePlayers
	} else {
		nextIdx = (nextIdx - 1 + numActivePlayers) % numActivePlayers
	}

	return nextIdx
}

// NextPlayer returns the next player in the play order
func (s *State) NextPlayer() *player.Player {
	if len(s.ActivePlayers) <= 1 {
		return nil
	}

	nextIdx := s.NextPlayerIndex()
	nextID := s.ActivePlayers[nextIdx]

	for _, p := range s.Players {
		if p.ID == nextID {
			return p
		}
	}

	return nil
}

// AdvanceTurn moves to the next player's turn
func (s *State) AdvanceTurn() {
	if len(s.ActivePlayers) <= 1 {
		return
	}

	s.CurrentPlayerIndex = s.NextPlayerIndex()
}

// SkipNextPlayer skips the next player's turn (used for Ace)
func (s *State) SkipNextPlayer() {
	if len(s.ActivePlayers) <= 2 {
		// With 2 players, skipping next is equivalent to playing again
		return
	}

	// Advance twice
	s.AdvanceTurn()
	s.AdvanceTurn()
}

// FindPlayerByID returns the player with the given ID
func (s *State) FindPlayerByID(id string) *player.Player {
	for _, p := range s.Players {
		if p.ID == id {
			return p
		}
	}
	return nil
}

// IsPlayerActive checks if a player is still active in the game
func (s *State) IsPlayerActive(playerID string) bool {
	for _, id := range s.ActivePlayers {
		if id == playerID {
			return true
		}
	}
	return false
}

// RemovePlayerFromActive removes a player from the active players list
func (s *State) RemovePlayerFromActive(playerID string) {
	for i, id := range s.ActivePlayers {
		if id == playerID {
			// Remove this player from active players
			s.ActivePlayers = append(s.ActivePlayers[:i], s.ActivePlayers[i+1:]...)

			// If the removed player was before the current player, adjust the index
			if i < s.CurrentPlayerIndex {
				s.CurrentPlayerIndex--
			}

			// If the removed player was the current player, don't change the index
			// (the next player will be at the same index position)
			if i == s.CurrentPlayerIndex && s.CurrentPlayerIndex >= len(s.ActivePlayers) {
				s.CurrentPlayerIndex = 0
			}

			return
		}
	}
}

// PlayCard plays the specified card from the player's hand
func (s *State) PlayCard(playerID string, c card.Card) error {
	// Check if it's the player's turn
	if playerID != s.CurrentPlayerID() {
		return errors.New("not your turn")
	}

	// Verify that the turn is not locked
	if s.LockedTurn {
		return errors.New("turn is locked")
	}

	// Find the player
	p := s.FindPlayerByID(playerID)
	if p == nil {
		return errors.New("player not found")
	}

	// Check if the player has the card
	if !p.HasCard(c) {
		return errors.New("card not in hand")
	}

	// Check if the play is valid
	if !player.CanPlayCardOn(c, s.TopCard, s.InAttackChain) {
		return errors.New("invalid play")
	}

	// If in an attack chain, only wild cards can be played on wild cards
	if s.InAttackChain && !c.IsWildCard() {
		return errors.New("must play a wild card to defend against an attack")
	}

	// Remove the card from the player's hand
	_, ok := p.RemoveFromHand(c)
	if !ok {
		return errors.New("failed to remove card from hand")
	}

	// Add the card to the discard pile
	s.DiscardPile = append(s.DiscardPile, c)
	s.TopCard = c

	// Update the last active suit (for Jack's suit change)
	if !c.IsJoker() {
		s.LastActiveSuit = c.Suit
	}

	// Handle wild cards
	if c.IsWildCard() {
		if s.InAttackChain {
			// Add to the attack amount
			s.AttackAmount += c.GetDrawPenalty()
		} else {
			// Start a new attack chain
			s.InAttackChain = true
			s.AttackAmount = c.GetDrawPenalty()
		}
	}

	// Check if the player has emptied their hand
	if p.HasEmptyHand() {
		s.RemovePlayerFromActive(playerID)
	}

	// Process special card effects
	s.ProcessCardEffect(c)

	return nil
}

// DrawCard makes the current player draw a card
func (s *State) DrawCard(playerID string) error {
	// Check if it's the player's turn
	if playerID != s.CurrentPlayerID() {
		return errors.New("not your turn")
	}

	// Find the player
	p := s.FindPlayerByID(playerID)
	if p == nil {
		return errors.New("player not found")
	}

	// Handle draw pile exhaustion
	if s.DrawPile.IsEmpty() {
		err := s.ReshuffleDiscardPile()
		if err != nil {
			return err
		}
	}

	// Draw a card for the player
	c, ok := s.DrawPile.Draw()
	if !ok {
		return errors.New("failed to draw card")
	}

	// Add the card to the player's hand
	p.AddToHand(c)

	// If in an attack chain, the player must draw the attack amount and end the chain
	if s.InAttackChain {
		// Draw the remaining attack amount - 1 (we already drew one)
		for i := 1; i < s.AttackAmount; i++ {
			// Handle draw pile exhaustion again if needed
			if s.DrawPile.IsEmpty() {
				err := s.ReshuffleDiscardPile()
				if err != nil {
					return err
				}
			}

			c, ok := s.DrawPile.Draw()
			if !ok {
				return errors.New("failed to draw attack penalty cards")
			}
			p.AddToHand(c)
		}

		// End the attack chain
		s.InAttackChain = false
		s.AttackAmount = 0
	}

	// Advance to the next player's turn
	s.AdvanceTurn()

	return nil
}

// ProcessCardEffect processes the effect of the played card
func (s *State) ProcessCardEffect(c card.Card) {
	if c.IsSkip() {
		// Ace skips the next player
		s.SkipNextPlayer()
	} else if c.IsWildCard() && s.InAttackChain {
		// In an attack chain, wild cards DO advance the turn
		// The next player must defend or draw
		s.AdvanceTurn()
	} else if c.IsSuitChanger() {
		// Jack changes the suit
		// So the turn is locked until the suit is changed
		s.LockTurn()
	} else if !s.InAttackChain {
		// In normal play, advance to the next player if not in an attack chain
		// and the card isn't a special card that changes turn order
		s.AdvanceTurn()
	}
}

// LockTurn locks the turn until the suit is changed
func (s *State) LockTurn() {
	s.LockedTurn = true
}

// UnlockTurn unlocks the turn
func (s *State) UnlockTurn() {
	s.LockedTurn = false
}

// ReshuffleDiscardPile reshuffles the discard pile (except top card) into the draw pile
func (s *State) ReshuffleDiscardPile() error {
	if len(s.DiscardPile) <= 1 {
		return errors.New("not enough cards to reshuffle")
	}

	// Keep the top card
	topCard := s.DiscardPile[len(s.DiscardPile)-1]

	// Add all other discard cards to draw pile
	cardsToShuffle := s.DiscardPile[:len(s.DiscardPile)-1]
	s.DrawPile.AddManyToBottom(cardsToShuffle)

	// Reset the discard pile with just the top card
	s.DiscardPile = []card.Card{topCard}

	// Shuffle the draw pile
	seed := rand.Int63()
	r := rand.New(rand.NewSource(seed))
	s.DrawPile.Shuffle(r)

	return nil
}

// ChangeSuit changes the active suit (for Jack effect)
func (s *State) ChangeSuit(playerID string, newSuit card.Suit) error {
	// Verify it's the player's turn
	if playerID != s.CurrentPlayerID() {
		return errors.New("not your turn")
	}

	// Verify that the turn is locked
	if !s.LockedTurn {
		return errors.New("turn is not locked")
	}

	// Verify that the last card played was a Jack
	if !s.TopCard.IsSuitChanger() {
		return errors.New("suit can only be changed after playing a Jack")
	}

	// Change the suit
	s.LastActiveSuit = newSuit
	s.UnlockTurn()
	s.AdvanceTurn()

	return nil
}

// IsGameOver checks if the game is over
func (s *State) IsGameOver() bool {
	return len(s.ActivePlayers) <= 1
}

// GetWinner returns the IDs of players who have won (emptied their hands)
func (s *State) GetWinner() []string {
	winners := make([]string, 0)

	// Anyone who emptied their hand is a winner
	for _, p := range s.Players {
		if p.HasEmptyHand() {
			winners = append(winners, p.ID)
		}
	}

	return winners
}

// GetLoser returns the ID of the last player left with cards
func (s *State) GetLoser() string {
	if len(s.ActivePlayers) != 1 {
		return ""
	}

	return s.ActivePlayers[0]
}

// ToJSON serializes the game state to JSON
func (s *State) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

// FromJSON deserializes the game state from JSON
func FromJSON(data []byte) (*State, error) {
	state := &State{}
	err := json.Unmarshal(data, state)
	return state, err
}
