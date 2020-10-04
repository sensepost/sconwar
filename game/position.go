package game

import "math/rand"

// Position is the position of an entity
type Position struct {
	X int
	Y int
}

// NewPosition returns a new position.
// The position itself is random
func NewPosition() *Position {

	x := rand.Intn(BoardX)
	y := rand.Intn(BoardY)

	return &Position{
		X: x,
		Y: y,
	}
}

// GetPosition returns the x, y of a position
func (p *Position) GetPosition() (int, int) {
	return p.X, p.Y
}

// MoveRandom moves to a random position for a
// specific distance
func (p *Position) MoveRandom(distance int) {

	if distance > BoardX || distance > BoardY {
		return
	}

	// X
	if rand.Float32() < 0.5 {
		p.X = p.X + distance
	} else {
		p.X = p.X - distance
	}

	// Y
	if rand.Float32() < 0.5 {
		p.Y = p.Y + distance
	} else {
		p.Y = p.Y - distance
	}

	p.floorAndCeilPosition()

}

// MoveTo moves to a specific x, y
func (p *Position) MoveTo(x int, y int) {

	p.X = x
	p.Y = y

	p.floorAndCeilPosition()
}

// floorAndCeilPosition ensures that the current x, y
// is in bounds.
func (p *Position) floorAndCeilPosition() {

	// x
	if p.X > BoardX {
		p.X = BoardX
	}

	if p.X < 0 {
		p.X = 0
	}

	if p.Y > BoardY {
		p.Y = BoardY
	}

	// y
	if p.Y < 0 {
		p.Y = 0
	}
}
