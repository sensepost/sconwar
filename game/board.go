package game

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
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
	mu            sync.RWMutex
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

// BoardInfoSnapshot is a read-only board view for API responses.
type BoardInfoSnapshot struct {
	Name          string
	Status        BoardStatus
	SizeX         int
	SizeY         int
	CurrentPlayer string
	FOWDistance   float64
	Created       time.Time
	Started       time.Time
	AliveCreep    int
	AlivePlayers  int
	PowerUps      int
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
			b.PowerUps = append(b.PowerUps, pup)
		}
	}

	gameState.With(prometheus.Labels{"state": "created"}).Inc()

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
	b.mu.Lock()
	b.Events = append(b.Events, event)
	b.mu.Unlock()

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
	b.mu.Lock()
	defer b.mu.Unlock()

	b.Players = append(b.Players, p)
}

// JoinPlayerIfOpen joins a player if this board accepts new players.
func (b *Board) JoinPlayerIfOpen(p *Player) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.Status != BoardStatusNew {
		return errors.New("game is not accepting new players")
	}

	for _, existing := range b.Players {
		if existing.ID == p.ID {
			return errors.New("player is already in the game")
		}
	}

	b.Players = append(b.Players, p)

	return nil
}

// FindPlayer gets a player by id.
func (b *Board) FindPlayer(playerID string) (*Player, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, p := range b.Players {
		if p.ID == playerID {
			return p, true
		}
	}

	return nil, false
}

// QueuePlayerAction queues an action for a board player.
func (b *Board) QueuePlayerAction(playerID string, action Action) error {
	player, ok := b.FindPlayer(playerID)
	if !ok {
		return errors.New("this player is not part of this game")
	}

	return player.AddAction(action)
}

// StatusValue gets the current board status.
func (b *Board) StatusValue() BoardStatus {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.Status
}

// SetStatus updates the board status.
func (b *Board) SetStatus(status BoardStatus) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Status = status
}

// SnapshotGameInfo returns immutable board info used in API responses.
func (b *Board) SnapshotGameInfo() BoardInfoSnapshot {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return BoardInfoSnapshot{
		Name:          b.Name,
		Status:        b.Status,
		SizeX:         b.SizeX,
		SizeY:         b.SizeY,
		CurrentPlayer: b.CurrentPlayer,
		FOWDistance:   b.FOWDistance,
		Created:       b.Created,
		Started:       b.Started,
		AliveCreep:    len(b.AliveCreep()),
		AlivePlayers:  len(b.AlivePlayers()),
		PowerUps:      len(b.PowerUps),
	}
}

// SnapshotEvents returns a copy of board events.
func (b *Board) SnapshotEvents() []*storage.Event {
	b.mu.RLock()
	defer b.mu.RUnlock()

	events := make([]*storage.Event, len(b.Events))
	copy(events, b.Events)
	return events
}

// SnapshotPlayers returns a copy of player pointers for read-only use.
func (b *Board) SnapshotPlayers() []*Player {
	b.mu.RLock()
	defer b.mu.RUnlock()

	players := make([]*Player, len(b.Players))
	copy(players, b.Players)
	return players
}

// SnapshotPowerUps returns a copy of powerup pointers.
func (b *Board) SnapshotPowerUps() []*PowerUp {
	b.mu.RLock()
	defer b.mu.RUnlock()

	powerups := make([]*PowerUp, len(b.PowerUps))
	copy(powerups, b.PowerUps)
	return powerups
}

// SnapshotCreeps returns a copy of creep pointers.
func (b *Board) SnapshotCreeps() []*Creep {
	b.mu.RLock()
	defer b.mu.RUnlock()

	creeps := make([]*Creep, len(b.Creeps))
	copy(creeps, b.Creeps)
	return creeps
}

// PlayerHasPowerup checks ownership of a powerup by a player.
func (b *Board) PlayerHasPowerup(playerID string, powerupID string) bool {
	player, ok := b.FindPlayer(playerID)
	if !ok {
		return false
	}

	for _, u := range player.PowerUps {
		if u.ID == powerupID {
			return true
		}
	}

	return false
}

// Run runs the game loop
func (b *Board) Run() {

	b.mu.Lock()
	b.Started = time.Now()
	b.mu.Unlock()
	b.SetStatus(BoardStatusRunning)
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
			b.mu.RLock()
			powerupCount := len(b.PowerUps)
			b.mu.RUnlock()
			if powerupCount < PowerUpMax {
				if pup := NewPowerUp(); pup != nil {
					b.mu.Lock()
					b.PowerUps = append(b.PowerUps, pup)
					b.mu.Unlock()
					b.LogEvent(&storage.Event{
						Date:        time.Now(),
						SrcEntity:   int(PowerupEntity),
					SrcEntityID: pup.ID,
					Msg:         `a new powerup has spawned on the board`,
				})
			}
		}

		b.processCreepTurn()

		for _, p := range b.AlivePlayers() {
			b.mu.Lock()
			b.CurrentPlayer = fmt.Sprintf(`(player) %s`, p.Name)
			b.mu.Unlock()
			ctx, cancel := context.WithTimeout(context.Background(), MaxRoundSeconds*time.Second)
			b.processPlayerTurn(ctx, p)
			cancel()
		}

		// win / stop conditions

		if b.StatusValue() != BoardStatusRunning {
			b.LogEvent(&storage.Event{
				Date: time.Now(),
				Msg: `game was in Run() loop, but status is not BoardStatusRunning. ` +
					`status may have changed from the outside (eg. via api)`,
			})

			// cleanup player scores
			for _, p := range b.AlivePlayers() {
				p.SaveFinalScore(b, b.CurrentDeathPosition())
			}

			break
		}

		if len(b.AlivePlayers()) == 1 && len(b.AliveCreep()) == 0 {
			player := b.AlivePlayers()[0]

			b.LogEvent(&storage.Event{
				Date:        time.Now(),
				SrcEntity:   int(PlayerEntity),
				SrcEntityID: player.ID,
				Msg:         fmt.Sprintf(`player %s is the last person standing`, player.Name),
			})

			player.SaveFinalScore(b, b.CurrentDeathPosition())
			gameState.With(prometheus.Labels{"state": "player-win"}).Inc()
			break
		}

		if len(b.AlivePlayers()) == 0 && len(b.AliveCreep()) > 0 {
			b.LogEvent(&storage.Event{
				Date: time.Now(),
				Msg:  `all players have been eliminated. the creep win!`,
			})
			gameState.With(prometheus.Labels{"state": "creep-win"}).Inc()
			break
		}
	}

	b.dbModel.Ended = time.Now()
	b.SetStatus(BoardStatusFinished)
	b.updateDbModel()
	gameState.With(prometheus.Labels{"state": "finished"}).Inc()
}

// RemovePowerUp removes a powerup from the board
func (b *Board) RemovePowerUp(powerup *PowerUp) {
	b.mu.Lock()
	defer b.mu.Unlock()

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
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, c := range b.Creeps {
		if c.Health > 0 {
			a = append(a, c)
		}
	}

	return
}

// AlivePlayers returns players that are alive
func (b *Board) AlivePlayers() (a []*Player) {
	b.mu.RLock()
	defer b.mu.RUnlock()
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
		b.mu.Lock()
		b.CurrentPlayer = fmt.Sprintf(`(creep) %s`, creep.Name)
		b.mu.Unlock()

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

			case Move:
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

						target.SaveFinalScore(b, b.CurrentDeathPosition())
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
			}
		}
	}
}

// chooseCreepAction randomly decides on an action to take
// with a bias towards moving and attacking
func (b *Board) chooseCreepAction() ActionType {

	c, err := wr.NewChooser(
		wr.Choice{Item: Move, Weight: 5},
		wr.Choice{Item: Attack, Weight: 4},
		wr.Choice{Item: Nothing, Weight: 2},
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to build creep action chooser")
		return Move
	}

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
