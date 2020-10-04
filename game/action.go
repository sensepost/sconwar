package game

import (
	"github.com/rs/zerolog/log"
)

// ActionType is an Action that can be invoked
type ActionType int

// Supported Actions
const (
	Move ActionType = iota
	Attack
	Pickup
)

// Action is the action that can be taken by a player
type Action struct {
	Action ActionType
	X      int
	Y      int
}

// NewAction starts a new action instance
func NewAction(action ActionType) *Action {
	return &Action{
		Action: action,
	}
}

// SetXY sets the x, t of an action
func (a *Action) SetXY(x int, y int) {
	a.X = x
	a.Y = y
}

// Execute executes an action
func (a *Action) Execute(player *Player) {
	switch a.Action {
	case Move:
		log.Info().Str("action", "move").Msg("executing a move command")
		player.MoveTo(a.X, a.Y)
		break
	}
}
