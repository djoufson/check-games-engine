package deck_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/deck"
)

// TestShouldContain54Cards_WhenCreatingNewDeck tests that a new deck has the correct card count
func TestShouldContain54Cards_WhenCreatingNewDeck(t *testing.T) {
	// Arrange & Act
	d := deck.New()

	// Assert
	if d.Count() != 54 {
		t.Errorf("Expected deck to have 54 cards, got %d", d.Count())
	}
}

// TestShouldContainCorrectSuitDistribution_WhenCreatingNewDeck tests suit distribution in a new deck
func TestShouldContainCorrectSuitDistribution_WhenCreatingNewDeck(t *testing.T) {
	// Arrange & Act
	d := deck.New()

	// Assert - Count cards by suit
	suitCounts := make(map[card.Suit]int)
	for _, c := range d.Cards {
		suitCounts[c.Suit]++
	}

	if suitCounts[card.Spades] != 13 {
		t.Errorf("Expected 13 spades, got %d", suitCounts[card.Spades])
	}
	if suitCounts[card.Hearts] != 13 {
		t.Errorf("Expected 13 hearts, got %d", suitCounts[card.Hearts])
	}
	if suitCounts[card.Diamonds] != 13 {
		t.Errorf("Expected 13 diamonds, got %d", suitCounts[card.Diamonds])
	}
	if suitCounts[card.Clubs] != 13 {
		t.Errorf("Expected 13 clubs, got %d", suitCounts[card.Clubs])
	}
	if suitCounts[card.Joker] != 2 {
		t.Errorf("Expected 2 jokers, got %d", suitCounts[card.Joker])
	}
}

// TestShouldContainBothColoredJokers_WhenCreatingNewDeck tests that both red and black jokers are present
func TestShouldContainBothColoredJokers_WhenCreatingNewDeck(t *testing.T) {
	// Arrange & Act
	d := deck.New()

	// Assert - Check for both joker colors
	jokerCount := 0
	redJokerFound := false
	blackJokerFound := false

	for _, c := range d.Cards {
		if c.Suit == card.Joker {
			jokerCount++
			if c.Color == card.Red {
				redJokerFound = true
			} else if c.Color == card.Black {
				blackJokerFound = true
			}
		}
	}

	if jokerCount != 2 {
		t.Errorf("Expected 2 jokers, got %d", jokerCount)
	}
	if !redJokerFound {
		t.Error("No red joker found in deck")
	}
	if !blackJokerFound {
		t.Error("No black joker found in deck")
	}
}
