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
	Player Player
	Action Actions
	X      int
	Y      int
}

func NewAction(player Player, action Actions) *Action {
	return &Action{
		Player: player,
		Action: action,
	}
}

func (a *Action) SetXY(x int, y int) {
	a.X = x
	a.Y = y
}

func (a *Action) Execute() {
	switch a.Action {
	case Move:
		log.Info().Str("action", "move").Msg("executing a move command")
		a.Player.MoveTo(a.X, a.Y)
		break
	}
}
