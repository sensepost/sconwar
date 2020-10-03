package game

import (
	"github.com/sensepost/sconwar/storage"
)

type Player struct {
	Name     string
	ID       string
	Health   uint
	Position *Position
	Command  int
}

func NewPlayer(p *storage.Player) *Player {
	return &Player{
		Name:     p.Name,
		ID:       p.UUID,
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

func (c *Player) DistanceFrom(o hasPosition) float64 {
	return distanceBetween(o, c.Position)
}
