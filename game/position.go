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

// NewPositionFromValue returns a new Position struct
// using the packed x & y value from v
func NewPositionFromValue(v int) *Position {

	x := v >> 32
	y := v - ((v >> 32) << 32)

	p := &Position{
		X: x, Y: y,
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

	// todo: add a chance that either x or y does not
	// change, meaning we just go up, down, left or
	// right. right now we only go diag.

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

// ToSingleValue packs the X, Y into a single value
func (p *Position) ToSingleValue() int {
	return (p.X << 32) + p.Y
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
