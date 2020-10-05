package game

import (
	"math/rand"

	"github.com/google/uuid"
)

// PowerUpType is the type of power up
type PowerUpType int

// Powerups
const (
	Health PowerUpType = iota
	Teleport
	DoubleDamage
)

// PowerUp is a powerup that a player could get
type PowerUp struct {
	ID       string      `json:"id"`
	Type     PowerUpType `json:"type"`
	Position *Position   `json:"position"`
}

// NewPowerUp creates a new random powerup.
// This is a chance based function and can return nil.
func NewPowerUp() *PowerUp {

	if rand.Intn(100) > PowerUpChance {
		return nil
	}

	// todo: add a chance that the powerup wont spawn,
	// say by setting an Invalid property when

	pups := []PowerUpType{Health, Teleport, DoubleDamage}

	return &PowerUp{
		ID:       uuid.New().String(),
		Type:     pups[rand.Intn(len(pups))],
		Position: NewPosition(),
	}
}

// GetPosition gets the players position
func (p *PowerUp) GetPosition() (int, int) {
	return p.Position.GetPosition()
}
