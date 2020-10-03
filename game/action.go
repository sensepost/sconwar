package game

import (
	"github.com/rs/zerolog/log"
)

type Actions int

const (
	Move Actions = iota
	Attack
	Pickup
)

type Action struct {
	Action Actions
	X      int
	Y      int
}

func NewAction(action Actions) *Action {
	return &Action{
		Action: action,
	}
}

func (a *Action) SetXY(x int, y int) {
	a.X = x
	a.Y = y
}

func (a *Action) Execute(player *Player) {
	switch a.Action {
	case Move:
		log.Info().Str("action", "move").Msg("executing a move command")
		player.MoveTo(a.X, a.Y)
		break
	}
}
