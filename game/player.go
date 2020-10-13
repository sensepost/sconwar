package game

import (
	"errors"
	"math/rand"

	"github.com/sensepost/sconwar/storage"
)

// ActionChannel is the channel a player will recieve actions on.
// This is actually a workaround as detailed in:
//	https://github.com/swaggo/swag/issues/680#issue-602785690
type ActionChannel chan Action

// Player represents a player of the game
type Player struct {
	Name   string `json:"name"`
	ID     string `json:"id"`
	Health int    `json:"health"`

	Position     *Position     `json:"position"`
	PowerUps     []*PowerUp    `json:"powerups"`
	PowerUpBuffs []PowerUpType `json:"buffs"`
	Actions      ActionChannel `json:"-"`
	ActionCount  int           `json:"action_count"`
}

// NewPlayer starts a new Player instance
func NewPlayer(p *storage.Player) *Player {

	return &Player{
		Name:     p.Name,
		ID:       p.UUID,
		Health:   100,
		Position: NewPosition(),
		Actions:  make(chan Action, RoundMoves),
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

	return nil
}

// Move moves the player to a random position
func (p *Player) Move() {
	p.Position.MoveRandom(1)
}

// MoveTo moves the player to a specific x.y
func (p *Player) MoveTo(x int, y int) {
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

// TakeDamage deals damage to the creep.
// An argument of -1 will make the damage taken
// random with a ceil of 30.
// The multiplier can be used to apply multiplication to the
// final damange taken.
func (p *Player) TakeDamage(dmg int, multiplier int) (int, int) {

	if dmg > 100 {
		dmg = -1
	}

	if dmg == -1 {
		dmg = rand.Intn(30)
	}

	dmg = dmg * multiplier

	p.Health -= dmg

	if p.Health < 0 {
		p.Health = 0
	}

	return dmg, p.Health
}

// GivePowerUp adds a powerup to the player
func (p *Player) GivePowerUp(powerup PowerUp) {
	p.PowerUps = append(p.PowerUps, &powerup)
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
	// todo: log events
	case Health:
		p.Health += PowerUpHealthBonus
		// todo: upper limit health to say 120?
		break
	case Teleport:
		p.PowerUpBuffs = append(p.PowerUpBuffs, Teleport)
		break
	case DoubleDamage:
		p.PowerUpBuffs = append(p.PowerUpBuffs, DoubleDamage)
		break
	}

	p.RemovePowerUp(powerup)
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
