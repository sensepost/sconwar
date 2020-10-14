package game

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/sensepost/sconwar/storage"

	wr "github.com/mroth/weightedrand"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// BoardStatus is the board status
type BoardStatus int

// The board statuses
const (
	BoardStatusNew BoardStatus = iota
	BoardStatusRunning
	BoardStatusFinished
)

// Board is the game board
type Board struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	Status        BoardStatus `json:"status"`
	SizeX         int         `json:"size_x"`
	SizeY         int         `json:"size_y"`
	FOWDistance   float64     `json:"fow_distance"`
	Events        []*storage.Event
	CurrentPlayer string `json:"current_player"`

	dbModel *storage.Board

	Creeps   []*Creep   `json:"creeps"`
	Players  []*Player  `json:"players"`
	PowerUps []*PowerUp `json:"powerups"`

	Created time.Time `json:"created"`
	Started time.Time `json:"started"`
}

// NewBoard starts a new Board
func NewBoard(id string, name string) *Board {

	b := &Board{
		ID:      id,
		Name:    name,
		Status:  BoardStatusNew,
		SizeX:   BoardX,
		SizeY:   BoardY,
		Created: time.Now(),
	}

	b.setFowDistance()

	// persist the board
	b.dbModel = &storage.Board{
		UUID:    id,
		Name:    b.Name,
		Created: b.Created,
	}
	storage.Storage.Get().Create(b.dbModel)

	for i := 0; i <= CreepCount; i++ {
		b.Creeps = append(b.Creeps, NewCreep())
	}

	for i := 0; i <= PowerUpMax; i++ {
		if pup := NewPowerUp(); pup != nil {
			// todo: log a new powerup spawning
			b.PowerUps = append(b.PowerUps, pup)
		}
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

func (b *Board) updateDbModel() *gorm.DB {
	return storage.Storage.Get().Save(b.dbModel)
}

// LogEvent logs a game event for this board
func (b *Board) LogEvent(event *storage.Event) {

	b.Events = append(b.Events, event)

	event.BoardID = b.dbModel.ID
	storage.Storage.Get().Create(&event)

	log.Info().
		Str("time", event.Date.Format(time.Stamp)).
		Str("src.entityid", event.SrcEntityID).
		Int("src.pos", event.SrcPos).
		Str("dst.entityid", event.DstEntityID).
		Int("dst.pos", event.DstPos).
		Str("msg", event.Msg).
		Msg("event")
}

// JoinPlayer joins a new human player to the board
func (b *Board) JoinPlayer(p *Player) {
	b.Players = append(b.Players, p)
}

// Run runs the game loop
func (b *Board) Run() {

	b.Started = time.Now()
	b.Status = BoardStatusRunning
	b.updateDbModel()

	for {

		log.Info().
			Str("board.id", b.ID).
			Int("creep.count", len(b.AliveCreep())).
			Int("player.count", len(b.AlivePlayers())).
			Msg("game stats")

		// respawn powerups. this is a chance-based thing so everytime
		// we get below the max amount of powerups we will try and
		// spawn new ones.
		if len(b.PowerUps) < PowerUpMax {
			if pup := NewPowerUp(); pup != nil {
				b.PowerUps = append(b.PowerUps, pup)
			}
		}

		b.processCreepTurn()

		for _, p := range b.AlivePlayers() {
			b.CurrentPlayer = p.ID
			ctx, cancel := context.WithTimeout(context.Background(), MaxRoundSeconds*time.Second)
			defer cancel()

			b.processPlayerTurn(ctx, p)
		}

		// win conditions
		if len(b.AlivePlayers()) == 1 && len(b.AliveCreep()) == 0 {
			player := b.AlivePlayers()[0]

			b.LogEvent(&storage.Event{
				Date:        time.Now(),
				SrcEntity:   int(PlayerEntity),
				SrcEntityID: player.ID,
				// todo: add player name
				Msg: fmt.Sprintf(`player %s is the last person standing`, player.Name),
			})

			player.SaveFinalScore(b.ID, b.CurrentDeathPosition())

			break
		}

		if len(b.AlivePlayers()) == 0 && len(b.AliveCreep()) > 0 {
			b.LogEvent(&storage.Event{
				Date: time.Now(),
				Msg:  `all players have been eliminated. the creep win!`,
			})
			break
		}
	}

	b.dbModel.Ended = time.Now()
	b.Status = BoardStatusFinished
	b.updateDbModel()
}

// RemovePowerUp removes a powerup from the board
func (b *Board) RemovePowerUp(powerup *PowerUp) {

	// copy the slice to be trimmed
	s := b.PowerUps

	for i, p := range b.PowerUps {
		if p != powerup {
			continue
		}

		s[i] = s[len(s)-1]
		b.PowerUps = s[:len(s)-1]

		break
	}

}

// AliveCreep returns creep that are alive
func (b *Board) AliveCreep() (a []*Creep) {
	for _, c := range b.Creeps {
		if c.Health > 0 {
			a = append(a, c)
		}
	}

	return
}

// AlivePlayers returns players that are alive
func (b *Board) AlivePlayers() (a []*Player) {
	for _, p := range b.Players {
		if p.Health > 0 {
			a = append(a, p)
		}
	}

	return
}

// TotalAliveEntities returns the # of alive entities in the game
func (b *Board) TotalAliveEntities() int {
	return len(b.AliveCreep()) + len(b.AlivePlayers())
}

// CurrentDeathPosition determines the current position assuming death
func (b *Board) CurrentDeathPosition() int {
	return b.TotalAliveEntities() + 1
}

// processCreepTurn processed the turn for each alive creep
// actions are chosen between moving and attacking.
// each creep will perform <CreepRoundMoves> number of moves.
func (b *Board) processCreepTurn() {

	for _, creep := range b.AliveCreep() {
		b.CurrentPlayer = creep.ID

		remMoves := CreepRoundMoves

		for remMoves > 0 {

			time.Sleep(time.Millisecond * 100) //todo: remove

			switch b.chooseCreepAction() {
			case Nothing:
				b.LogEvent(&storage.Event{
					Date:        time.Now(),
					SrcEntity:   int(CreepEntity),
					SrcEntityID: creep.ID,
					Action:      int(Nothing),
					Msg:         fmt.Sprintf(`creep %s decided to do nothing`, creep.Name),
				})

				remMoves--
				break

			case Move:
				// todo: move events are pretty noisy, maybe we don't need to record those?
				sourcepos := creep.Position.ToSingleValue()

				creep.Move()

				b.LogEvent(&storage.Event{
					Date:        time.Now(),
					SrcEntity:   int(CreepEntity),
					SrcEntityID: creep.ID,
					SrcPos:      sourcepos,
					DstPos:      creep.Position.ToSingleValue(),
					Action:      int(Move),
					Msg:         fmt.Sprintf(`creep %s moved position`, creep.Name),
				})

				remMoves--
				break

			case Attack:
				// process alive players before alive creep
				for _, target := range b.AlivePlayers() {
					if !creep.IsInAttackRangeOf(target) {
						continue
					}

					dmg, h := target.TakeDamage(-1, 1)

					b.LogEvent(&storage.Event{
						Date:        time.Now(),
						SrcEntity:   int(CreepEntity),
						SrcEntityID: creep.ID,
						SrcPos:      creep.Position.ToSingleValue(),
						DstEntity:   int(PlayerEntity),
						DstEntityID: target.ID,
						DstPos:      target.Position.ToSingleValue(),
						Action:      int(Attack),
						Msg: fmt.Sprintf(`creep %s attacked player %s for %d damage`,
							creep.Name, target.Name, dmg),
					})

					if h == 0 {
						b.LogEvent(&storage.Event{
							Date:        time.Now(),
							SrcEntity:   int(CreepEntity),
							SrcEntityID: target.ID,
							SrcPos:      target.Position.ToSingleValue(),
							Action:      int(Attack),
							Msg:         fmt.Sprintf(`player %s has been killed`, target.Name),
						})

						target.SaveFinalScore(b.ID, b.CurrentDeathPosition())
					}

					remMoves--

					break // alivePlayers loop
				}

				for _, target := range b.AliveCreep() {
					// prevent creep suicide
					if creep.ID == target.ID {
						continue
					}

					if !creep.IsInAttackRangeOf(target) {
						continue
					}

					dmg, h := target.TakeDamage(-1, 1)

					b.LogEvent(&storage.Event{
						Date:        time.Now(),
						SrcEntity:   int(CreepEntity),
						SrcEntityID: creep.ID,
						SrcPos:      creep.Position.ToSingleValue(),
						DstEntity:   int(CreepEntity),
						DstEntityID: target.ID,
						DstPos:      target.Position.ToSingleValue(),
						Action:      int(Attack),
						Msg: fmt.Sprintf(`creep %s attacked creep %s for %d damage`,
							creep.Name, target.Name, dmg),
					})

					if h == 0 {
						b.LogEvent(&storage.Event{
							Date:        time.Now(),
							SrcEntity:   int(CreepEntity),
							SrcEntityID: target.ID,
							SrcPos:      target.Position.ToSingleValue(),
							Action:      int(Attack),
							Msg:         fmt.Sprintf(`creep %s has been killed`, target.Name),
						})
					}

					remMoves--

					break // aliveCreep loop
				}

				break // attack
			}
		}
	}
}

// chooseCreepAction randomly decides on an action to take
// with a boas towards moving and attacking
func (b *Board) chooseCreepAction() ActionType {

	c := wr.NewChooser(
		wr.Choice{Item: Move, Weight: 5},
		wr.Choice{Item: Attack, Weight: 5},
		wr.Choice{Item: Nothing, Weight: 2},
	)

	return c.Pick().(ActionType)
}

// processPlayerTurn executes the actions for a player
func (b *Board) processPlayerTurn(ctx context.Context, p *Player) {

	// only execute the max number of actions
	for i := 0; i < PlayerRoundMoves; i++ {
		select {
		case <-ctx.Done():
			b.LogEvent(&storage.Event{
				Date:        time.Now(),
				SrcEntity:   int(PlayerEntity),
				SrcEntityID: p.ID,
				SrcPos:      p.Position.ToSingleValue(),
				Msg:         fmt.Sprintf(`timeout while waiting for player %s`, p.Name),
			})
			return
		case action := <-p.Actions:
			action.Execute(p, b)
			p.ActionCount--
		}
	}
}
