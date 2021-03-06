package game

import (
	"math/rand"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
)

// Creep is a creep
type Creep struct {
	Name     string    `json:"name"`
	Position *Position `json:"position"`
	Health   int       `json:"health"`
	ID       string    `json:"id"`
}

// NewCreep returns a new initialised Creep
func NewCreep() *Creep {
	return &Creep{
		Name:     strings.ToLower(randomdata.SillyName()),
		Position: NewPosition(),
		ID:       uuid.New().String(),
		Health:   100,
	}
}

// Move moves a creep in a random position
func (c *Creep) Move() {
	c.Position.MoveRandom(MaxCreepMoveDistance)
	distanceMovedByEntity.With(prometheus.Labels{"entity": "creep"}).Inc()
}

// GetPosition gets the x, y position of a creep
func (c *Creep) GetPosition() (int, int) {
	return c.Position.GetPosition()
}

// IsInAttackRangeOf checks if something is within attack range
func (c *Creep) IsInAttackRangeOf(o hasPosition) bool {
	distance := distanceBetween(o, c.Position)

	if distance <= AttackRange {
		return true
	}

	return false
}

// TakeDamage deals damage to the creep.
// An argument of -1 will make the damage taken
// random with a ceil of 30.
// The multiplier can be used to apply multiplication to the
// final damange taken.
func (c *Creep) TakeDamage(dmg int, multiplier int) (int, int) {

	if dmg > 100 {
		dmg = -1
	}

	if dmg == -1 {
		dmg = rand.Intn(MaxDamage)
	}

	dmg = dmg * multiplier

	c.Health -= dmg

	if c.Health < 0 {
		c.Health = 0
	}

	damageTaken.With(prometheus.Labels{"entity": "creep"}).Add(float64(dmg))
	return dmg, c.Health
}
