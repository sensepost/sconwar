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
		p.X = p.X - distance
	}

	if p.X > BoardX {
		p.X = BoardX
	}

	if p.X < 0 {
		p.X = 0
	}

	// Y
	if rand.Float32() < 0.5 {
		p.Y = p.Y + distance
	} else {
		p.Y = p.Y - distance
	}

	if p.Y > BoardY {
		p.Y = BoardY
	}

	if p.Y < 0 {
		p.Y = 0
	}
}
