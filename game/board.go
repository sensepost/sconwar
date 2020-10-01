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

		time.Sleep(1 * time.Second)
		log.Info().
			Str("board.id", b.ID).
			Int("creep.count", len(b.Creeps)).
			Int("player.count", len(b.Players)).
			Msg("10 sec sleep")

		// cleanup dead creep/people
	}
}

func (b *Board) moveAndAttackCreep() {

	for _, creep := range b.Creeps {

		creep.Move()

		for _, target := range b.Creeps {
			if creep == target {
				continue
			}

			if creep.IsInRangeOf(target) {

				dmg := target.TakeDamage(-1)
				log.Warn().
					Str("game", b.ID).
					Str("attacker", creep.ID).
					Str("victim", target.ID).
					Int("damage", dmg).
					Msg("attacked in range creep")
			}
		}
	}
}
