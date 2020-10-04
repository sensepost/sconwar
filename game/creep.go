package game

import (
	"math/rand"

	"github.com/google/uuid"
)

// Creep is a creep
type Creep struct {
	Position *Position
	Health   int
	ID       string
}

// NewCreep returns a new initialised Creep
func NewCreep() *Creep {
	// todo: choose random names for creep
	return &Creep{
		Position: NewPosition(),
		ID:       uuid.New().String(),
		Health:   100,
	}
}

// Move moves a creep in a random position
func (c *Creep) Move() {
	c.Position.MoveRandom(1)
}

// GetPosition gets the x, y position of a creep
func (c *Creep) GetPosition() (int, int) {
	return c.Position.GetPosition()
}

// IsInRangeOf checks if something is within attack range
func (c *Creep) IsInRangeOf(o hasPosition) bool {
	distance := distanceBetween(o, c.Position)

	if distance <= AttackRange {
		return true
	}

	return false
}

// TakeDamage deals damage to the creep.
// An argument of -1 will make the damage taken
// random with a ceil of 30
func (c *Creep) TakeDamage(dmg int) (int, int) {

	if dmg > 100 {
		dmg = -1
	}

	if dmg == -1 {
		dmg = rand.Intn(30)
	}

	c.Health -= dmg

	if c.Health < 0 {
		c.Health = 0
	}

	return dmg, c.Health
}
