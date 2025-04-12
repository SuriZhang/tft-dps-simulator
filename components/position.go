package components

// Position represents a champion's position on the board
type Position struct {
    X float64
    Y float64
}

// NewPosition creates a Position component
func NewPosition(x, y float64) Position {
    return Position{
        X: x,
        Y: y,
    }
}