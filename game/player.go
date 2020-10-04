package game

import (
	"errors"

	"github.com/sensepost/sconwar/storage"
)

// Player represents a player of the game
type Player struct {
	Name     string
	ID       string
	Health   uint
	Position *Position
	Actions  chan Action `json:"-" swaggerignore:"true"`
}

// NewPlayer starts a new Player instance
func NewPlayer(p *storage.Player) *Player {

	return &Player{
		Name:     p.Name,
		ID:       p.UUID,
		Health:   100,
		Position: NewPosition(),
		Actions:  make(chan Action, RoundMoves),
	}
}

// AddAction queues a new action for a player
func (p *Player) AddAction(action Action) error {
	select {
	case p.Actions <- action:
	default:
		return errors.New(`player action buffer full`)
	}

	return nil
}

// Move moves the player to a random position
func (p *Player) Move() {
	p.Position.MoveRandom(1)
}

// MoveTo moves the player to a specific x.y
func (p *Player) MoveTo(x int, y int) {
	p.Position.MoveTo(x, y)
}

// GetPosition gets the players position
func (p *Player) GetPosition() (int, int) {
	return p.Position.GetPosition()
}

// DistanceFrom calculates the distance from another entity
func (p *Player) DistanceFrom(o hasPosition) float64 {
	return distanceBetween(o, p.Position)
}
