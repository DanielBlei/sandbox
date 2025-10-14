# Knights Game

A Go console game challenge where six knights battle until only one remains standing.

*This is a challenge solution focused on core game logic.

## Game Rules

- **Six knights** start the game (K1, K2, K3, K4, K5, K6)
- Each knight begins with **10 health points**
- Moves are made in **clockwise order**
- Every knight hits the next one for a **random damage** (1-6 health points)
- When a knight's health reaches 0, they **die and are excluded** from the game
- The **last surviving knight wins**

## How to Play

```bash
cd knights_game
go run main.go
```

### Reproducible Games

For testing or to replay the same game sequence, you can specify a seed:

```bash
go run main.go -seed 12345
```

## Example Output

```
K1 hits K2 for 4
K2 hits K3 for 2
K3 hits K4 for 6
K4 hits K5 for 3
K5 hits K6 for 5
K6 hits K1 for 1
...
K3 wins
```

## Implementation Details

- Uses seeded random number generation for damage calculation (testable and reproducible)
- Implements proper turn management and knight elimination
- Handles edge cases like skipping dead knights in turn order
- Command-line seed support for deterministic gameplay
