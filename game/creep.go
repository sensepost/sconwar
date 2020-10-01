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

func NewCreep() *Creep {

	return &Creep{
		Position: NewPosition(),
		ID:       uuid.New().String(),
		Health:   100,
	}
}

func (c *Creep) Move() {
	c.Position.MoveRandom(1)
}

func (p *Creep) GetPosition() (int, int) {
	return p.Position.GetPosition()
}

func (c *Creep) IsInRangeOf(o hasPosition) bool {
	distance := distanceBetween(o, c.Position)

	if distance <= 2 {
		return true
	}

	return false
}

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
