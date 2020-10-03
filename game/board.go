package game

import (
	"context"
	"math"
	"time"

	"github.com/rs/zerolog/log"
)

// Board is the game board
type Board struct {
	ID          string
	SizeX       int
	SizeY       int
	FOWDistance float64
	Creeps      []*Creep
	Players     []*Player

	Created time.Time
	Started time.Time
}

// NewBoard starts a new Board
func NewBoard(id string) *Board {

	b := &Board{
		ID:      id,
		SizeX:   BoardX,
		SizeY:   BoardY,
		Created: time.Now(),
	}

	b.setFowDistance()

	for i := 0; i <= CreepCount; i++ {
		b.Creeps = append(b.Creeps, NewCreep())
	}

	return b
}

// setFowDistance calculates the visible fog of war distance
func (b *Board) setFowDistance() {

	first := math.Pow(float64(BoardX), 2)
	second := math.Pow(float64(BoardY), 2)

	distance := math.Sqrt(first + second)

	b.FOWDistance = distance / 100 * FogOfWarPercent
}

// JoinPlayer joins a new human player to the board
func (b *Board) JoinPlayer(p *Player) {
	b.Players = append(b.Players, p)
}

// Run runs the game loop
func (b *Board) Run() {

	b.Started = time.Now()

	for {

		log.Info().
			Str("board.id", b.ID).
			Int("creep.count", len(b.aliveCreep())).
			Int("player.count", len(b.alivePlayers())).
			Msg("game stats")

		b.moveAndAttackCreep()

		for _, p := range b.alivePlayers() {
			ctx, cancel := context.WithTimeout(context.Background(), MaxRoundSeconds*time.Second)
			defer cancel()

			b.processPlayerActions(ctx, p)
		}

		// todo: cleanup dead creep/people

		if len(b.aliveCreep()) == 1 {
			log.Error().Msg("Game finished, last man standing!")
			return
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

func (b *Board) alivePlayers() (a []*Player) {

	for _, p := range b.Players {
		if p.Health > 0 {
			a = append(a, p)
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

				// todo: limit moves. atm the creep will
				// shoot all others in range

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

// processPlayerActions executes the actions for a player
func (b *Board) processPlayerActions(ctx context.Context, p *Player) {

	// only execute the max number of actions
	for i := 0; i < RoundMoves; i++ {
		select {
		case <-ctx.Done():
			return
		case action := <-p.Actions:
			action.Execute(p)
		}
	}
}
