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
	c.Position.MoveTo(10, 10)
	oldX, oldY := c.Position.X, c.Position.Y
	c.Move()
	if c.Position.X == oldX && c.Position.Y == oldY {
		t.Fatalf("expected Cords to change but remained identical")
	}
}
