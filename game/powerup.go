package game

import (
	"math/rand"

	"github.com/google/uuid"
	wr "github.com/mroth/weightedrand"
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

	c := wr.NewChooser(
		wr.Choice{Item: Health, Weight: 5},
		wr.Choice{Item: Teleport, Weight: 5},
		wr.Choice{Item: DoubleDamage, Weight: 5},
	)

	return &PowerUp{
		ID:       uuid.New().String(),
		Type:     c.Pick().(PowerUpType),
		Position: NewPosition(),
	}
}

// GetPosition gets the players position
func (p *PowerUp) GetPosition() (int, int) {
	return p.Position.GetPosition()
}
