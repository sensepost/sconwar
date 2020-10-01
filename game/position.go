package game

import "math/rand"

// Position is the position of an entity
type Position struct {
	X int
	Y int
}

func NewPosition() *Position {

	x := rand.Intn(BoardX)
	y := rand.Intn(BoardY)

	return &Position{
		X: x,
		Y: y,
	}
}

func (p *Position) GetPosition() (int, int) {
	return p.X, p.Y
}

func (p *Position) MoveRandom(distance int) {

	if distance > BoardX || distance > BoardY {
		return
	}

	// X
	if rand.Float32() < 0.5 {
		p.X = p.X + distance
	} else {
		// TODO, bounds check
		p.X = p.X - distance
	}

	// Y
	if rand.Float32() < 0.5 {
		p.Y = p.Y + distance
	} else {
		// TODO, bounds check
		p.Y = p.Y - distance
	}
}
