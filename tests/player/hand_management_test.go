package player_test

import (
	"testing"

	"github.com/djoufson/check-games-engine/card"
	"github.com/djoufson/check-games-engine/player"
)

// TestShouldIncreaseHandSize_WhenAddingSingleCard tests adding a single card to a player's hand
func TestShouldIncreaseHandSize_WhenAddingSingleCard(t *testing.T) {
	// Arrange
	p := player.New("player1")
	c := card.NewCard(card.Spades, card.Ace)

	// Act
	p.AddToHand(c)

	// Assert
	if len(p.Hand) != 1 {
		t.Errorf("Expected hand to have 1 card, got %d", len(p.Hand))
	}
}

// TestShouldAddCorrectCardToHand_WhenAddingSingleCard tests the correct card is added to hand
func TestShouldAddCorrectCardToHand_WhenAddingSingleCard(t *testing.T) {
	// Arrange
	p := player.New("player1")
	c := card.NewCard(card.Spades, card.Ace)

	// Act
	p.AddToHand(c)

	// Assert
	if p.Hand[0].Suit != card.Spades || p.Hand[0].Rank != card.Ace {
		t.Errorf("Expected Ace of Spades, got %v", p.Hand[0])
	}
}

// TestShouldIncreaseHandByCorrectAmount_WhenAddingMultipleCards tests adding multiple cards
func TestShouldIncreaseHandByCorrectAmount_WhenAddingMultipleCards(t *testing.T) {
	// Arrange
	p := player.New("player1")
	cards := []card.Card{
		card.NewCard(card.Spades, card.Ace),
		card.NewCard(card.Hearts, card.King),
	}

	// Act
	p.AddCardsToHand(cards)

	// Assert
	if len(p.Hand) != 2 {
		t.Errorf("Expected hand to have 2 cards, got %d", len(p.Hand))
	}
}

// TestShouldRemoveCardSuccessfully_WhenCardExistsInHand tests successful card removal
func TestShouldRemoveCardSuccessfully_WhenCardExistsInHand(t *testing.T) {
	// Arrange
	p := player.New("player1")
	aceSpades := card.NewCard(card.Spades, card.Ace)
	kingHearts := card.NewCard(card.Hearts, card.King)
	p.AddCardsToHand([]card.Card{aceSpades, kingHearts})

	// Act
	removedCard, ok := p.RemoveFromHand(aceSpades)

	// Assert
	if !ok {
		t.Error("Expected card to be removed successfully")
	}
	if removedCard.Suit != card.Spades || removedCard.Rank != card.Ace {
		t.Errorf("Expected to remove Ace of Spades, got %v", removedCard)
	}
}

// TestShouldDecreaseHandSize_WhenRemovingCard tests hand size decreases after removal
func TestShouldDecreaseHandSize_WhenRemovingCard(t *testing.T) {
	// Arrange
	p := player.New("player1")
	aceSpades := card.NewCard(card.Spades, card.Ace)
	kingHearts := card.NewCard(card.Hearts, card.King)
	p.AddCardsToHand([]card.Card{aceSpades, kingHearts})

	// Act
	_, _ = p.RemoveFromHand(aceSpades)

	// Assert
	if len(p.Hand) != 1 {
		t.Errorf("Expected hand to have 1 card after removal, got %d", len(p.Hand))
	}
}

// TestShouldReturnFailure_WhenRemovingNonexistentCard tests removing a card that doesn't exist
func TestShouldReturnFailure_WhenRemovingNonexistentCard(t *testing.T) {
	// Arrange
	p := player.New("player1")
	aceSpades := card.NewCard(card.Spades, card.Ace)
	kingHearts := card.NewCard(card.Hearts, card.King)
	p.AddCardsToHand([]card.Card{kingHearts})

	// Act
	_, ok := p.RemoveFromHand(aceSpades)

	// Assert
	if ok {
		t.Error("Expected removal of non-existent card to fail")
	}
}

// TestShouldDetectCard_WhenCardIsInHand tests checking if a card is in the hand
func TestShouldDetectCard_WhenCardIsInHand(t *testing.T) {
	// Arrange
	p := player.New("player1")
	aceSpades := card.NewCard(card.Spades, card.Ace)
	p.AddCardsToHand([]card.Card{aceSpades})

	// Act & Assert
	if !p.HasCard(aceSpades) {
		t.Error("Expected player to have Ace of Spades")
	}
}

// TestShouldNotDetectCard_WhenCardIsNotInHand tests checking if a card is not in the hand
func TestShouldNotDetectCard_WhenCardIsNotInHand(t *testing.T) {
	// Arrange
	p := player.New("player1")
	aceSpades := card.NewCard(card.Spades, card.Ace)
	kingHearts := card.NewCard(card.Hearts, card.King)
	p.AddCardsToHand([]card.Card{aceSpades})

	// Act & Assert
	if p.HasCard(kingHearts) {
		t.Error("Expected player not to have King of Hearts")
	}
}
