package game

import (
	"errors"

	"github.com/sensepost/sconwar/storage"
)

type Player struct {
	Name     string
	ID       string
	Health   uint
	Position *Position
	Actions  chan Action `json:"-" swaggerignore:"true"`
}

func NewPlayer(p *storage.Player) *Player {

	return &Player{
		Name:     p.Name,
		ID:       p.UUID,
		Health:   100,
		Position: NewPosition(),
		Actions:  make(chan Action, RoundMoves),
	}
}

func (p *Player) AddAction(action Action) error {
	select {
	case p.Actions <- action:
	default:
		return errors.New(`player action buffer full`)
	}

	return nil
}

func (p *Player) Move() {
	p.Position.MoveRandom(1)
}

func (p *Player) MoveTo(x int, y int) {
	p.Position.MoveTo(x, y)
}

func (p *Player) GetPosition() (int, int) {
	return p.Position.GetPosition()
}

func (c *Player) DistanceFrom(o hasPosition) float64 {
	return distanceBetween(o, c.Position)
}
