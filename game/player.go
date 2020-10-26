package game

import (
	"errors"
	"math/rand"

	"github.com/dariubs/percent"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sensepost/sconwar/storage"
)

// ActionChannel is the channel a player will recieve actions on.
// This is actually a workaround as detailed in:
//	https://github.com/swaggo/swag/issues/680#issue-602785690
type ActionChannel chan Action

// Player represents a player of the game
type Player struct {
	Name   string `json:"name"`
	ID     string `json:"-"` // don't leak the id in the api
	Health int    `json:"health"`

	Position     *Position     `json:"position"`
	PowerUps     []*PowerUp    `json:"powerups"`
	PowerUpBuffs []PowerUpType `json:"buffs"`
	Actions      ActionChannel `json:"-"`
	ActionCount  int           `json:"action_count"`

	Score         int `json:"score"`
	DamageTaken   int `json:"damage_taken"`
	DamageDealt   int `json:"damage_dealt"`
	CreepKilled   int `json:"killed_creep"`
	PlayersKilled int `json:"killed_players"`
}

// NewPlayer starts a new Player instance
func NewPlayer(p *storage.Player) *Player {

	playerCreated.WithLabelValues().Inc()
	return &Player{
		Name:        p.Name,
		ID:          p.UUID,
		Health:      100,
		Position:    NewPosition(),
		Actions:     make(chan Action, PlayerRoundMoves),
		Score:       0,
		DamageDealt: 0,
		DamageTaken: 0,
	}
}

// AddAction queues a new action for a player
func (p *Player) AddAction(action Action) error {
	select {
	case p.Actions <- action:
	default:
		return errors.New(`player action buffer full`)
	}

	p.ActionCount++
	playerActionsQueued.WithLabelValues().Inc()
	return nil
}

// Move moves the player to a random position
func (p *Player) Move() {
	distanceMovedByEntity.With(prometheus.Labels{"entity": "player"}).Inc()
	p.Position.MoveRandom(1)
}

// MoveTo moves the player to a specific x.y
func (p *Player) MoveTo(x int, y int) {
	distanceMovedByEntity.With(prometheus.Labels{"entity": "player"}).Inc()
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

// IsInAttackRangeOf checks if something is within attack range
func (p *Player) IsInAttackRangeOf(o hasPosition) bool {
	distance := distanceBetween(o, p.Position)

	if distance <= AttackRange {
		return true
	}

	return false
}

// TakeDamage deals damage to the player.
// An argument of -1 will make the damage taken
// random with a ceil of 30.
// The multiplier can be used to apply multiplication to the
// final damange taken.
func (p *Player) TakeDamage(dmg int, multiplier int) (int, int) {

	if dmg > 100 {
		dmg = -1
	}

	if dmg == -1 {
		dmg = rand.Intn(MaxDamage)
	}

	dmg = dmg * multiplier

	p.Health -= dmg

	if p.Health < 0 {
		p.Health = 0
	}

	p.AddDamageTaken(dmg)
	damageTaken.With(prometheus.Labels{"entity": "player"}).Add(float64(dmg))
	return dmg, p.Health
}

// GivePowerUp adds a powerup to the player
func (p *Player) GivePowerUp(powerup PowerUp) {
	p.PowerUps = append(p.PowerUps, &powerup)
	powerupsCollected.WithLabelValues().Inc()
}

// HasAvailableBuf checks of a player has a buff enabled
func (p *Player) HasAvailableBuf(buf PowerUpType) bool {
	for _, u := range p.PowerUpBuffs {
		if u == buf {
			return true
		}
	}

	return false
}

// RemoveBuf removes a buf from a player
func (p *Player) RemoveBuf(buf PowerUpType) {
	bufs := p.PowerUpBuffs

	for i, u := range bufs {
		if u == buf {
			bufs[i] = bufs[len(bufs)-1]
			p.PowerUpBuffs = bufs[:len(bufs)-1]

			break
		}
	}
}

// UsePowerUp uses a powerup, applying the relevant buf
func (p *Player) UsePowerUp(powerupID string) {

	var powerup *PowerUp
	for _, u := range p.PowerUps {
		if powerupID == u.ID {
			powerup = u
		}
	}

	// this shouldn't happen, but ok just in case
	if powerup == nil {
		return
	}

	switch powerup.Type {
	case Health:
		p.Health += PowerUpHealthBonus
		// ciel max health
		if p.Health > PowerUpHealthBonusMax {
			p.Health = PowerUpHealthBonusMax
		}
		break
	case Teleport:
		p.PowerUpBuffs = append(p.PowerUpBuffs, Teleport)
		break
	case DoubleDamage:
		p.PowerUpBuffs = append(p.PowerUpBuffs, DoubleDamage)
		break
	}

	p.RemovePowerUp(powerup)
	powerupsUsed.With(prometheus.Labels{"poweruptype": string(powerup.Type)}).Inc()
}

// RemovePowerUp removes a powerup from the player
func (p *Player) RemovePowerUp(powerup *PowerUp) {

	// copy the slice to be trimmed
	s := p.PowerUps

	for i, u := range p.PowerUps {
		if u != powerup {
			continue
		}

		s[i] = s[len(s)-1]
		p.PowerUps = s[:len(s)-1]

		break
	}
}

// AddScore adds to the players existing score
func (p *Player) AddScore(amount int) {
	p.Score += amount
	scoreAwarded.WithLabelValues().Add(float64(amount))
}

// AddDamageTaken adds to the players existing damage taken
func (p *Player) AddDamageTaken(amount int) {
	p.DamageTaken += amount
}

// AddDamageDealt adds to the players existing damage dealt
func (p *Player) AddDamageDealt(amount int) {
	p.DamageDealt += amount
}

// RecordCreepKilled records a score for a killed creep
func (p *Player) RecordCreepKilled() {
	p.Score += CreepKilledScore
	p.CreepKilled++
	creepsKilled.WithLabelValues().Inc()
}

// RecordPlayerKilled records a score for a killed creep
func (p *Player) RecordPlayerKilled() {
	p.Score += PlayerKilledScore
	p.PlayersKilled++
	playersKilled.WithLabelValues().Inc()
}

// SaveFinalScore stores the players score
func (p *Player) SaveFinalScore(board *Board, position int) {

	positionPercent := percent.PercentOf(position, (len(board.Players) + len(board.Creeps)))

	potentialBonus := 0
	potentialBonus += len(board.Players) * PlayerBonusScore
	potentialBonus += len(board.Creeps) * CreepBonusScore

	// bonus is the remainder of the position % as a % of the potential bonus
	bonus := int(percent.Percent(int(100-positionPercent), potentialBonus))

	var player storage.Player
	storage.Storage.Get().Where("uuid = ?", p.ID).First(&player)

	score := &storage.PlayerGameScore{
		PlayerID:      player.ID,
		BoardID:       board.ID,
		Position:      position,
		Score:         p.Score + bonus,
		DamageTaken:   p.DamageTaken,
		DamageDealt:   p.DamageDealt,
		CreepKilled:   p.CreepKilled,
		PlayersKilled: p.PlayersKilled,
	}
	storage.Storage.Get().Create(&score)
}
