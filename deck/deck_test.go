package deck

import (
	"math/rand"
	"testing"
	"time"

	"github.com/djoufson/check-games-engine/card"
)

func TestNewDeck(t *testing.T) {
	deck := New()

	// A standard deck should have 54 cards (52 normal + 2 jokers)
	if deck.Count() != 54 {
		t.Errorf("Expected deck to have 54 cards, got %d", deck.Count())
	}

	// Check that there are 13 cards of each suit
	suitCounts := make(map[card.Suit]int)
	for _, c := range deck.Cards {
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

	// Check jokers have appropriate colors
	jokerCount := 0
	redJokerFound := false
	blackJokerFound := false

	for _, c := range deck.Cards {
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

func TestShuffle(t *testing.T) {
	deck1 := New()
	deck2 := New()

	// Before shuffling, the decks should be identical
	for i := range deck1.Cards {
		if deck1.Cards[i].Suit != deck2.Cards[i].Suit || deck1.Cards[i].Rank != deck2.Cards[i].Rank {
			t.Errorf("Decks are not identical before shuffling at index %d", i)
		}
	}

	// Shuffle one deck
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	deck1.Shuffle(r)

	// Decks should now be different (this could theoretically fail, but probability is extremely low)
	identical := true
	for i := range deck1.Cards {
		if deck1.Cards[i].Suit != deck2.Cards[i].Suit || deck1.Cards[i].Rank != deck2.Cards[i].Rank {
			identical = false
			break
		}
	}

	if identical {
		t.Error("Decks are still identical after shuffling")
	}
}

func TestDraw(t *testing.T) {
	deck := New()
	initialCount := deck.Count()

	// Draw a card
	_, ok := deck.Draw()
	if !ok {
		t.Error("Failed to draw a card from a non-empty deck")
	}

	// Deck should have one less card
	if deck.Count() != initialCount-1 {
		t.Errorf("Expected deck to have %d cards after drawing, got %d", initialCount-1, deck.Count())
	}

	// Drawing from an empty deck
	emptyDeck := New()
	emptyDeck.Cards = nil
	_, ok = emptyDeck.Draw()
	if ok {
		t.Error("Expected drawing from empty deck to return false")
	}
}

func TestDrawN(t *testing.T) {
	deck := New()
	initialCount := deck.Count()

	// Draw 5 cards
	cards, ok := deck.DrawN(5)
	if !ok {
		t.Error("Failed to draw 5 cards from a full deck")
	}

	if len(cards) != 5 {
		t.Errorf("Expected to draw 5 cards, got %d", len(cards))
	}

	// Deck should have 5 fewer cards
	if deck.Count() != initialCount-5 {
		t.Errorf("Expected deck to have %d cards after drawing 5, got %d", initialCount-5, deck.Count())
	}

	// Try to draw more cards than available
	_, ok = deck.DrawN(100)
	if ok {
		t.Error("Expected drawing more cards than available to return false")
	}
}

func TestAddCards(t *testing.T) {
	deck := New()
	initialCount := deck.Count()

	// Draw a card to add back later
	drawnCard, _ := deck.Draw()

	// Add a card to the bottom
	deck.AddToBottom(drawnCard)
	if deck.Count() != initialCount {
		t.Errorf("Expected deck to have %d cards after adding to bottom, got %d", initialCount, deck.Count())
	}

	// The added card should be the last one
	if deck.Cards[deck.Count()-1].Suit != drawnCard.Suit || deck.Cards[deck.Count()-1].Rank != drawnCard.Rank {
		t.Error("Card added to bottom is not the last card in the deck")
	}

	// Draw again to test adding to top
	drawnCard, _ = deck.Draw()

	// Add a card to the top
	deck.AddToTop(drawnCard)

	// The added card should be the first one
	if deck.Cards[0].Suit != drawnCard.Suit || deck.Cards[0].Rank != drawnCard.Rank {
		t.Error("Card added to top is not the first card in the deck")
	}
}

func TestAddManyToBottom(t *testing.T) {
	deck := New()

	// Draw 10 cards
	drawnCards, _ := deck.DrawN(10)
	initialCount := deck.Count()

	// Add the cards back to the bottom
	deck.AddManyToBottom(drawnCards)

	if deck.Count() != initialCount+10 {
		t.Errorf("Expected deck to have %d cards after adding many to bottom, got %d", initialCount+10, deck.Count())
	}
}

func TestIsEmpty(t *testing.T) {
	deck := New()
	if deck.IsEmpty() {
		t.Error("New deck should not be empty")
	}

	// Draw all cards
	deck.DrawN(deck.Count())

	if !deck.IsEmpty() {
		t.Error("Deck should be empty after drawing all cards")
	}
}
