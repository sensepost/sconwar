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
func (a *Action) Execute(player *Player, board *Board) {
	switch a.Action {
	case Move:
		log.Info().Str("action", "move").Msg("executing a move command")
		player.MoveTo(a.X, a.Y)
		break
	case Attack:
		log.Info().Str("action", "attack").Msg("executing an attack command")

		// find entities on the x, y and if there is something, take damage
		for _, c := range board.aliveCreep() {
			cx, cy := c.GetPosition()
			if cx == a.X && cy == a.Y {
				dmg, h := c.TakeDamage(-1)
				log.Info().Int("damage", dmg).Int("health", h).Msg("attacked creep")
			}
		}

		for _, p := range board.alivePlayers() {
			px, py := p.GetPosition()
			if px == a.X && py == a.Y {
				dmg, h := p.TakeDamage(-1)
				log.Info().Int("damage", dmg).Int("health", h).Msg("attacked player")
			}
		}

		// todo: log this in the player/creep structs.
		// todo: resolve names for the creep / player that is attacked
	}
}
