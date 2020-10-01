package game

import (
	"time"

	"github.com/rs/zerolog/log"
)

// Board is the game board
type Board struct {
	ID      string
	SizeX   int
	SizeY   int
	Creeps  []*Creep
	Players []*Player
}

// NewBoard starts a new Board
func NewBoard(id string) *Board {

	b := &Board{
		ID:    id,
		SizeX: BoardX,
		SizeY: BoardY,
	}

	for i := 0; i <= CreepCount; i++ {
		b.Creeps = append(b.Creeps, NewCreep())
	}

	return b
}

func (b *Board) Run() {

	for {

		b.moveAndAttackCreep()

		time.Sleep(10 * time.Millisecond)
		log.Info().
			Str("board.id", b.ID).
			Int("creep.count", len(b.aliveCreep())).
			Int("player.count", len(b.Players)).
			Msg("10 sec sleep")

		// cleanup dead creep/people
		if len(b.aliveCreep()) == 1 {
			log.Error().Msg("Game finished, last man standing!")
		}
	}
}

func (b *Board) aliveCreep() (a []*Creep) {

	for _, c := range b.Creeps {
		if c.Health > 0 {
			a = append(a, c)
		}
	}

	return
}

func (b *Board) moveAndAttackCreep() {

	for _, creep := range b.aliveCreep() {

		creep.Move()

		for _, target := range b.aliveCreep() {
			if creep == target {
				continue
			}

			if creep.IsInRangeOf(target) {

				dmg, h := target.TakeDamage(-1)
				log.Warn().
					Str("game", b.ID).
					Str("attacker", creep.ID).
					Str("victim", target.ID).
					Int("damage", dmg).
					Int("health", h).
					Msg("attacked in range creep")

				if h == 0 {
					log.Error().Str("target.id", target.ID).Msg("creep has been killed!")
				}
			}
		}
	}
}
