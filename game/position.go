package game

import "math/rand"

// Position is the position of an entity
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// NewPosition returns a new position.
// The position itself is random
func NewPosition() *Position {

	p := &Position{
		X: rand.Intn(BoardX),
		Y: rand.Intn(BoardY),
	}

	p.floorAndCeilPosition()

	return p
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

	if p.X < 1 {
		p.X = 1
	}

	if p.Y > BoardY {
		p.Y = BoardY
	}

	// y
	if p.Y < 1 {
		p.Y = 1
	}
}
