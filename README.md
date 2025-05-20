# Check-Game Engine

A Go implementation of the check-game rules engine, designed to be portable, testable, and embeddable in various platforms.

## Overview

This engine implements the card game logic as defined in the rules.md file, focusing on:

- Core game mechanics
- Move validation
- Special card effects
- Game state management
- No UI or persistence layer

## Project Structure

```
game-engine/
├── card/          # Card data types and utilities
├── game/          # Game logic implementation
├── deck/          # Deck management and shuffling
├── player/        # Player state and actions
├── state/         # Game state and transitions
├── validation/    # Move validation logic
└── tests/         # Test cases
```

## Features

- Stateless pure functions for game logic
- Extensive test coverage for all rules
- Serializable game state (JSON)
- Support for all special cards and their effects:
  - Aces (skip next player)
  - 7s (draw 2 cards)
  - Jokers (draw 4 cards)
  - Jacks (change suit)
  - 2s (transparent/wildcard)
- Attack chain handling
- Deterministic for testing

## Usage

### Basic Usage

```go
import "github.com/check-games/engine/game"

// Create a new game with 3 players
gameState, err := game.NewGame([]string{"player1", "player2", "player3"}, nil)

// Play a card
newState, err := game.PlayCard(gameState, "player1", card.Card{Suit: card.Spades, Rank: card.King})

// Check if move is valid
isValid, err := game.ValidateMove(gameState, "player1", card.Card{Suit: card.Spades, Rank: card.King})

// Draw a card
newState, err := game.DrawCard(gameState, "player2")
```

## Testing

Run the tests with:

```bash
go test ./...
```

## License

MIT 