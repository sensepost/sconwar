package game

import "github.com/google/uuid"

type Player struct {
	Name     string
	ID       string
	Health   uint
	Position *Position
}

func NewPlayer(name string) *Player {
	return &Player{
		Name:     name,
		ID:       uuid.New().String(),
		Health:   100,
		Position: NewPosition(),
	}
}

func (p *Player) Move() {
	p.Position.MoveRandom(1)
}

func (p *Player) GetPosition() (int, int) {
	return p.Position.GetPosition()
}
