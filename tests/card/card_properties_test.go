package card_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
)

// TestShouldIdentifySkipCards_WhenCheckingCardProperties tests if Ace cards are properly identified as skip cards
func TestShouldIdentifySkipCards_WhenCheckingCardProperties(t *testing.T) {
	// Arrange
	aceCard := card.NewCard(card.Spades, card.Ace)
	normalCard := card.NewCard(card.Hearts, card.Three)

	// Act & Assert
	if !aceCard.IsSkip() {
		t.Errorf("Expected Ace to be a skip card")
	}

	if normalCard.IsSkip() {
		t.Errorf("Expected Three not to be a skip card")
	}
}

// TestShouldIdentifyWildCards_WhenCheckingCardProperties tests if Seven and Joker cards are properly identified as wild cards
func TestShouldIdentifyWildCards_WhenCheckingCardProperties(t *testing.T) {
	// Arrange
	sevenCard := card.NewCard(card.Hearts, card.Seven)
	jokerCard := card.NewRedJoker()
	normalCard := card.NewCard(card.Clubs, card.Five)

	// Act & Assert
	if !sevenCard.IsWildCard() {
		t.Errorf("Expected Seven to be a wild card")
	}

	if !jokerCard.IsWildCard() {
		t.Errorf("Expected Joker to be a wild card")
	}

	if normalCard.IsWildCard() {
		t.Errorf("Expected Five not to be a wild card")
	}
}

// TestShouldIdentifySuitChangers_WhenCheckingCardProperties tests if Jack cards are properly identified as suit changers
func TestShouldIdentifySuitChangers_WhenCheckingCardProperties(t *testing.T) {
	// Arrange
	jackCard := card.NewCard(card.Diamonds, card.Jack)
	normalCard := card.NewCard(card.Clubs, card.Nine)

	// Act & Assert
	if !jackCard.IsSuitChanger() {
		t.Errorf("Expected Jack to be a suit changer")
	}

	if normalCard.IsSuitChanger() {
		t.Errorf("Expected Nine not to be a suit changer")
	}
}

// TestShouldIdentifyTransparentCards_WhenCheckingCardProperties tests if Two cards are properly identified as transparent
func TestShouldIdentifyTransparentCards_WhenCheckingCardProperties(t *testing.T) {
	// Arrange
	twoCard := card.NewCard(card.Clubs, card.Two)
	normalCard := card.NewCard(card.Diamonds, card.Eight)

	// Act & Assert
	if !twoCard.IsTransparent() {
		t.Errorf("Expected Two to be a transparent card")
	}

	if normalCard.IsTransparent() {
		t.Errorf("Expected Eight not to be a transparent card")
	}
}

// TestShouldIdentifyJokers_WhenCheckingCardProperties tests if Joker cards are properly identified
func TestShouldIdentifyJokers_WhenCheckingCardProperties(t *testing.T) {
	// Arrange
	redJoker := card.NewRedJoker()
	blackJoker := card.NewBlackJoker()
	normalCard := card.NewCard(card.Spades, card.King)

	// Act & Assert
	if !redJoker.IsJoker() {
		t.Errorf("Expected Red Joker to be identified as a joker")
	}

	if !blackJoker.IsJoker() {
		t.Errorf("Expected Black Joker to be identified as a joker")
	}

	if normalCard.IsJoker() {
		t.Errorf("Expected King not to be identified as a joker")
	}
}

// TestShouldCalculateCorrectDrawPenalty_WhenCheckingWildCards tests the draw penalty values
func TestShouldCalculateCorrectDrawPenalty_WhenCheckingWildCards(t *testing.T) {
	// Arrange
	sevenCard := card.NewCard(card.Hearts, card.Seven)
	jokerCard := card.NewRedJoker()
	normalCard := card.NewCard(card.Clubs, card.Six)

	// Act & Assert
	if sevenCard.GetDrawPenalty() != 2 {
		t.Errorf("Expected Seven to have draw penalty of 2, got %d", sevenCard.GetDrawPenalty())
	}

	if jokerCard.GetDrawPenalty() != 4 {
		t.Errorf("Expected Joker to have draw penalty of 4, got %d", jokerCard.GetDrawPenalty())
	}

	if normalCard.GetDrawPenalty() != 0 {
		t.Errorf("Expected normal card to have draw penalty of 0, got %d", normalCard.GetDrawPenalty())
	}
}
