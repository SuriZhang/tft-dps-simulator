package components

// Position represents a champion's position on the board
type Position struct {
	x int
	y int
}

// NewPosition creates a Position component
func NewPosition(x, y int) Position {
	return Position{
		x: x,
		y: y,
	}
}

// GetX returns the x-coordinate of the position
func (p Position) GetX() int {
	return p.x
}

// GetY returns the y-coordinate of the position
func (p Position) GetY() int {
	return p.y
}

// SetX sets the x-coordinate of the position
func (p *Position) SetX(x int) {
	p.x = x
}

// SetY sets the y-coordinate of the position
func (p *Position) SetY(y int) {
	p.y = y
}

// SetPosition sets both x and y coordinates of the position
func (p *Position) SetPosition(x, y int) {
	p.x = x
	p.y = y
}
