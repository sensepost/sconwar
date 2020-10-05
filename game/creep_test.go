package game_test

import (
	"testing"

	"github.com/sensepost/sconwar/game"
)

func TestCreep_CanCreateNew(t *testing.T) {
	c := game.NewCreep()
	if c.Health != 100 {
		t.Fatalf("expected 100, got: %d", c.Health)
	}
}

func TestCreep_CanMove(t *testing.T) {
	c := game.NewCreep()
	pos := c.Position
	c.Move()
	if c.Position.X == pos.X && c.Position.Y == pos.Y {
		t.Fatalf("expected Cords to change but remained identical")
	}
}
