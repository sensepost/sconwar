package game

import (
	"strings"
	"testing"

	"github.com/sensepost/sconwar/storage"
)

func mustCreateStoragePlayer(t *testing.T, name string, id string) *storage.Player {
	t.Helper()

	p := &storage.Player{
		Name: name,
		UUID: id,
	}
	storage.Storage.Get().Create(p)
	return p
}

func lastEventMsg(t *testing.T, b *Board) string {
	t.Helper()

	events := b.SnapshotEvents()
	if len(events) == 0 {
		t.Fatalf("expected at least one event")
	}

	return events[len(events)-1].Msg
}

func anyEventContains(b *Board, needle string) bool {
	for _, e := range b.SnapshotEvents() {
		if strings.Contains(e.Msg, needle) {
			return true
		}
	}
	return false
}

func TestRoundFeatures_PickupUseMoveAndTeleport(t *testing.T) {
	resetGameTestEnv(t)
	InitMetrics()

	b := NewBoard("round-1", "round-features")
	b.Creeps = nil
	b.PowerUps = nil

	sp := mustCreateStoragePlayer(t, "p1", "p1-id")
	p := NewPlayer(sp)
	p.Position.MoveTo(5, 5)
	b.JoinPlayer(p)

	healthPU := &PowerUp{ID: "pu-health", Type: Health, Position: &Position{X: 6, Y: 5}}
	teleportPU := &PowerUp{ID: "pu-tele", Type: Teleport, Position: &Position{X: 7, Y: 5}}
	b.PowerUps = []*PowerUp{healthPU, teleportPU}

	p.Health = 80

	pickHealth := NewAction(Pickup)
	pickHealth.SetXY(6, 5)
	pickHealth.Execute(p, b)
	if len(p.PowerUps) != 1 {
		t.Fatalf("expected 1 powerup after pickup, got %d", len(p.PowerUps))
	}

	p.UsePowerUp("pu-health")
	if p.Health <= 80 {
		t.Fatalf("expected health increase after using health powerup, got %d", p.Health)
	}
	if p.Health > PowerUpHealthBonusMax {
		t.Fatalf("expected health <= %d, got %d", PowerUpHealthBonusMax, p.Health)
	}

	pickTeleport := NewAction(Pickup)
	pickTeleport.SetXY(7, 5)
	pickTeleport.Execute(p, b)
	if len(p.PowerUps) != 1 {
		t.Fatalf("expected 1 powerup after teleport pickup, got %d", len(p.PowerUps))
	}

	p.UsePowerUp("pu-tele")
	if !p.HasAvailableBuf(Teleport) {
		t.Fatalf("expected teleport buff to be available")
	}

	moveFar := NewAction(Move)
	moveFar.SetXY(15, 15)
	moveFar.Execute(p, b)

	x, y := p.GetPosition()
	if x != 15 || y != 15 {
		t.Fatalf("expected teleport move to reach 15,15 got %d,%d", x, y)
	}
	if p.HasAvailableBuf(Teleport) {
		t.Fatalf("expected teleport buff to be consumed after far move")
	}
}

func TestRoundFeatures_InvalidAndOutOfRangeActions(t *testing.T) {
	resetGameTestEnv(t)
	InitMetrics()

	b := NewBoard("round-2", "invalid-actions")
	b.Creeps = nil
	b.PowerUps = nil

	sp := mustCreateStoragePlayer(t, "p1", "p2-id")
	p := NewPlayer(sp)
	p.Position.MoveTo(5, 5)
	b.JoinPlayer(p)

	outOfRangeMove := NewAction(Move)
	outOfRangeMove.SetXY(20, 20)
	outOfRangeMove.Execute(p, b)
	msg := lastEventMsg(t, b)
	if !strings.Contains(msg, "out of range") {
		t.Fatalf("expected out-of-range message, got: %s", msg)
	}
	x, y := p.GetPosition()
	if x != 5 || y != 5 {
		t.Fatalf("expected position unchanged for out-of-range move, got %d,%d", x, y)
	}

	invalid := NewAction(ActionType(99))
	invalid.Execute(p, b)
	if !anyEventContains(b, "invalid action") {
		t.Fatalf("expected invalid action message in event stream")
	}
}

func TestRoundFeatures_AttackAndPickupOutOfRangeNoOp(t *testing.T) {
	resetGameTestEnv(t)
	InitMetrics()

	b := NewBoard("round-3", "attack-pickup")
	b.PowerUps = nil
	b.Creeps = nil

	sp1 := mustCreateStoragePlayer(t, "p1", "pa-1")
	sp2 := mustCreateStoragePlayer(t, "p2", "pa-2")
	p1 := NewPlayer(sp1)
	p2 := NewPlayer(sp2)
	p1.Position.MoveTo(5, 5)
	p2.Position.MoveTo(6, 5)
	b.JoinPlayer(p1)
	b.JoinPlayer(p2)

	attack := NewAction(Attack)
	attack.SetXY(6, 5)
	attack.Execute(p1, b)
	msg := lastEventMsg(t, b)
	if !strings.Contains(msg, "attacked player") {
		t.Fatalf("expected attack event, got: %s", msg)
	}

	pu := &PowerUp{ID: "far-pu", Type: DoubleDamage, Position: &Position{X: 20, Y: 20}}
	b.PowerUps = []*PowerUp{pu}

	pickFar := NewAction(Pickup)
	pickFar.SetXY(20, 20)
	pickFar.Execute(p1, b)
	if len(p1.PowerUps) != 0 {
		t.Fatalf("expected no pickup for out-of-range powerup")
	}
}
