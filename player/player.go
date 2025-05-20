// Package player provides player-related functionality for the check-game engine.
package player

import (
	"github.com/check-games/engine/card"
)

// Player represents a player in the game
type Player struct {
	ID   string      `json:"id"`
	Hand []card.Card `json:"hand"`
}

// New creates a new player with the given ID and an empty hand
func New(id string) *Player {
	return &Player{
		ID:   id,
		Hand: make([]card.Card, 0),
	}
}

// AddToHand adds a card to the player's hand
func (p *Player) AddToHand(c card.Card) {
	p.Hand = append(p.Hand, c)
}

// AddCardsToHand adds multiple cards to the player's hand
func (p *Player) AddCardsToHand(cards []card.Card) {
	p.Hand = append(p.Hand, cards...)
}

// RemoveFromHand removes a card from the player's hand and returns it
// Returns false if the card is not in the hand
func (p *Player) RemoveFromHand(c card.Card) (card.Card, bool) {
	for i, handCard := range p.Hand {
		if handCard.Suit == c.Suit && handCard.Rank == c.Rank && handCard.Color == c.Color {
			// Remove the card from the hand
			p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)
			return c, true
		}
	}
	return card.Card{}, false
}

// HasCard returns true if the player has the specified card in their hand
func (p *Player) HasCard(c card.Card) bool {
	for _, handCard := range p.Hand {
		if handCard.Suit == c.Suit && handCard.Rank == c.Rank && handCard.Color == c.Color {
			return true
		}
	}
	return false
}

// HasMatchingCard returns true if the player has a card that matches the specified
// card by color, rank, or suit according to the game rules
func (p *Player) HasMatchingCard(c card.Card, includeWildCards bool) bool {
	for _, handCard := range p.Hand {
		// Transparent card (2) can be played on any card
		if handCard.IsTransparent() {
			return true
		}

		// Wild cards can be played on other wild cards
		if includeWildCards && handCard.IsWildCard() && c.IsWildCard() {
			return true
		}

		// Regular matching: same suit (color for jokers) or same rank
		if handCard.Suit == c.Suit {
			return true
		}

		if handCard.Rank == c.Rank {
			return true
		}

		// Joker color matching
		if handCard.IsJoker() && c.Color == handCard.Color {
			return true
		}
	}
	return false
}

// GetPlayableCards returns a list of cards that the player can play on the specified card
func (p *Player) GetPlayableCards(c card.Card, inAttackChain bool) []card.Card {
	playable := make([]card.Card, 0)

	for _, handCard := range p.Hand {
		if CanPlayCardOn(handCard, c, inAttackChain) {
			playable = append(playable, handCard)
		}
	}

	return playable
}

// CanPlayCardOn checks if a card can be played on another card according to the rules
func CanPlayCardOn(playedCard, topCard card.Card, inAttackChain bool) bool {
	// In an attack chain, only wild cards can be played on wild cards
	if inAttackChain {
		// If we're in an attack chain, only wild cards can respond to wild cards
		if topCard.IsWildCard() {
			return playedCard.IsWildCard()
		}
		return false
	}

	// Transparent card (2) and Suit changer (Jack) can be played on any card except during attack chain
	if playedCard.IsTransparent() || playedCard.IsSuitChanger() {
		return true
	}

	// Wild cards can be played on other wild cards (7 or Joker)
	if playedCard.IsWildCard() && topCard.IsWildCard() {
		return true
	}

	// Regular matching: same suit or same rank
	if playedCard.Suit == topCard.Suit {
		return true
	}

	if playedCard.Rank == topCard.Rank {
		return true
	}

	// Joker color matching
	if playedCard.IsJoker() && topCard.Color == playedCard.Color {
		return true
	}

	return false
}

// HandSize returns the number of cards in the player's hand
func (p *Player) HandSize() int {
	return len(p.Hand)
}

// HasEmptyHand returns true if the player has no cards left
func (p *Player) HasEmptyHand() bool {
	return len(p.Hand) == 0
}
