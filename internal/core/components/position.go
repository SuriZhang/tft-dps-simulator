package components

// Position represents a champion's position on the board
type Position struct {
	x float64
	y float64
}

// NewPosition creates a Position component
func NewPosition(x, y float64) Position {
	return Position{
		x: x,
		y: y,
	}
}

// GetX returns the x-coordinate of the position
func (p Position) GetX() float64 {
	return p.x
}

// GetY returns the y-coordinate of the position
func (p Position) GetY() float64 {
	return p.y
}

// SetX sets the x-coordinate of the position
func (p *Position) SetX(x float64) {
	p.x = x
}

// SetY sets the y-coordinate of the position
func (p *Position) SetY(y float64) {
	p.y = y
}

// SetPosition sets both x and y coordinates of the position
func (p *Position) SetPosition(x, y float64) {
	p.x = x
	p.y = y
}
