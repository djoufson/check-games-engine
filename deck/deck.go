// Package deck provides functionality for creating and managing a deck of cards.
package deck

import (
	"math/rand"

	"github.com/djoufson/check-games-engine/card"
)

// Deck represents a deck of cards
type Deck struct {
	Cards []card.Card
}

// New creates a new standard deck of cards (52 cards + 2 jokers)
func New() *Deck {
	cards := make([]card.Card, 0, 54)

	// Add all standard cards
	suits := []card.Suit{card.Spades, card.Hearts, card.Diamonds, card.Clubs}
	ranks := []card.Rank{
		card.Ace, card.Two, card.Three, card.Four, card.Five, card.Six, card.Seven,
		card.Eight, card.Nine, card.Ten, card.Jack, card.Queen, card.King,
	}

	for _, suit := range suits {
		for _, rank := range ranks {
			cards = append(cards, card.NewCard(suit, rank))
		}
	}

	// Add jokers
	cards = append(cards, card.NewRedJoker())
	cards = append(cards, card.NewBlackJoker())

	return &Deck{Cards: cards}
}

// Shuffle randomizes the order of cards in the deck
func (d *Deck) Shuffle(r *rand.Rand) {
	// Fisher-Yates shuffle
	for i := len(d.Cards) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}
}

// Draw removes and returns the top card from the deck
func (d *Deck) Draw() (card.Card, bool) {
	if len(d.Cards) == 0 {
		return card.Card{}, false
	}

	card := d.Cards[0]
	d.Cards = d.Cards[1:]
	return card, true
}

// DrawN draws n cards from the top of the deck
func (d *Deck) DrawN(n int) ([]card.Card, bool) {
	if len(d.Cards) < n {
		return nil, false
	}

	cards := d.Cards[:n]
	d.Cards = d.Cards[n:]
	return cards, true
}

// AddToBottom adds a card to the bottom of the deck
func (d *Deck) AddToBottom(c card.Card) {
	d.Cards = append(d.Cards, c)
}

// AddToTop adds a card to the top of the deck
func (d *Deck) AddToTop(c card.Card) {
	d.Cards = append([]card.Card{c}, d.Cards...)
}

// AddManyToBottom adds multiple cards to the bottom of the deck
func (d *Deck) AddManyToBottom(cards []card.Card) {
	d.Cards = append(d.Cards, cards...)
}

// Count returns the number of cards in the deck
func (d *Deck) Count() int {
	return len(d.Cards)
}

// IsEmpty returns true if the deck has no cards
func (d *Deck) IsEmpty() bool {
	return len(d.Cards) == 0
}
