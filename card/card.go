// Package card defines the card types and utilities for the check-game engine.
package card

import (
	"encoding/json"
	"fmt"
)

// Suit represents the suit of a playing card
type Suit string

// Card suits
const (
	Spades   Suit = "SPADES"
	Hearts   Suit = "HEARTS"
	Diamonds Suit = "DIAMONDS"
	Clubs    Suit = "CLUBS"
	Joker    Suit = "JOKER" // Special suit for jokers
)

// Rank represents the rank of a playing card
type Rank string

// Card ranks
const (
	Ace   Rank = "ACE"
	Two   Rank = "TWO"
	Three Rank = "THREE"
	Four  Rank = "FOUR"
	Five  Rank = "FIVE"
	Six   Rank = "SIX"
	Seven Rank = "SEVEN"
	Eight Rank = "EIGHT"
	Nine  Rank = "NINE"
	Ten   Rank = "TEN"
	Jack  Rank = "JACK"
	Queen Rank = "QUEEN"
	King  Rank = "KING"
)

// Color represents the color of a playing card
type Color string

// Card colors
const (
	Red   Color = "RED"
	Black Color = "BLACK"
)

// Card represents a playing card with a suit, rank, and color
type Card struct {
	Suit  Suit  `json:"suit"`
	Rank  Rank  `json:"rank"`
	Color Color `json:"color"`
}

// NewCard creates a new card with the given suit and rank, and automatically assigns the color
func NewCard(suit Suit, rank Rank) Card {
	var color Color

	switch suit {
	case Hearts, Diamonds:
		color = Red
	case Spades, Clubs:
		color = Black
	case Joker:
		// For Joker, we need to specify a color explicitly
		if rank == "RED" {
			color = Red
			rank = ""
		} else {
			color = Black
			rank = ""
		}
	}

	return Card{
		Suit:  suit,
		Rank:  rank,
		Color: color,
	}
}

// NewRedJoker creates a new red joker card
func NewRedJoker() Card {
	return Card{
		Suit:  Joker,
		Rank:  "",
		Color: Red,
	}
}

// NewBlackJoker creates a new black joker card
func NewBlackJoker() Card {
	return Card{
		Suit:  Joker,
		Rank:  "",
		Color: Black,
	}
}

// String returns a string representation of the card
func (c Card) String() string {
	if c.Suit == Joker {
		return fmt.Sprintf("%s Joker", c.Color)
	}
	return fmt.Sprintf("%s of %s", c.Rank, c.Suit)
}

// MarshalJSON provides custom JSON marshaling
func (c Card) MarshalJSON() ([]byte, error) {
	type Alias Card
	return json.Marshal(&struct {
		Alias
		Display string `json:"display"`
	}{
		Alias:   Alias(c),
		Display: c.String(),
	})
}

// IsJoker returns true if the card is a joker
func (c Card) IsJoker() bool {
	return c.Suit == Joker
}

// IsWildCard returns true if the card is a wild card (7 or Joker)
func (c Card) IsWildCard() bool {
	return c.Rank == Seven || c.IsJoker()
}

// IsTransparent returns true if the card is transparent (rank 2)
func (c Card) IsTransparent() bool {
	return c.Rank == Two
}

// IsSkip returns true if the card skips the next player (Ace)
func (c Card) IsSkip() bool {
	return c.Rank == Ace
}

// IsSuitChanger returns true if the card can change the suit (Jack)
func (c Card) IsSuitChanger() bool {
	return c.Rank == Jack
}

// GetDrawPenalty returns the number of cards to draw as penalty for this card
func (c Card) GetDrawPenalty() int {
	if c.Rank == Seven {
		return 2
	}
	if c.IsJoker() {
		return 4
	}
	return 0
} 