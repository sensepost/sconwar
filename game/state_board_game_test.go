package game

import (
	"path/filepath"
	"strconv"
	"sync"
	"testing"

	"github.com/sensepost/sconwar/storage"
)

func resetGameTestEnv(t *testing.T) {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "db.sqlite")
	if err := storage.InitDbPath(dbPath); err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	Setup()
}

func TestState_CreateGetListGames(t *testing.T) {
	Setup()

	CreateGame("g1", &Board{ID: "g1", Name: "one"})
	CreateGame("g2", &Board{ID: "g2", Name: "two"})

	g, ok := GetGame("g1")
	if !ok {
		t.Fatalf("expected game g1 to exist")
	}
	if g.Name != "one" {
		t.Fatalf("expected game name one, got %s", g.Name)
	}

	all := ListGames()
	if len(all) != 2 {
		t.Fatalf("expected 2 games, got %d", len(all))
	}
}

func TestState_ConcurrentCreate(t *testing.T) {
	Setup()

	const total = 100
	wg := sync.WaitGroup{}
	wg.Add(total)

	for i := 0; i < total; i++ {
		i := i
		go func() {
			defer wg.Done()
			id := "g" + strconv.Itoa(i)
			CreateGame(id, &Board{ID: id, Name: "n"})
			_, _ = GetGame(id)
		}()
	}

	wg.Wait()

	if len(ListGames()) != total {
		t.Fatalf("expected %d games, got %d", total, len(ListGames()))
	}
}

func TestBoard_JoinPlayerIfOpenAndDuplicateGuard(t *testing.T) {
	b := &Board{Status: BoardStatusNew}
	p := &Player{ID: "p1", Name: "player-one", Actions: make(chan Action, PlayerRoundMoves)}

	if err := b.JoinPlayerIfOpen(p); err != nil {
		t.Fatalf("unexpected join error: %v", err)
	}
	if err := b.JoinPlayerIfOpen(p); err == nil {
		t.Fatalf("expected duplicate player join to fail")
	}

	b.SetStatus(BoardStatusRunning)
	if err := b.JoinPlayerIfOpen(&Player{ID: "p2", Actions: make(chan Action, PlayerRoundMoves)}); err == nil {
		t.Fatalf("expected join on running board to fail")
	}
}

func TestBoard_QueuePlayerActionMembershipAndQueue(t *testing.T) {
	p := &Player{
		ID:      "p1",
		Actions: make(chan Action, PlayerRoundMoves),
	}
	b := &Board{
		Status:  BoardStatusRunning,
		Players: []*Player{p},
	}

	a := NewAction(Move)
	a.SetXY(1, 1)
	if err := b.QueuePlayerAction("p1", *a); err != nil {
		t.Fatalf("expected queue success, got: %v", err)
	}
	if p.ActionCount != 1 {
		t.Fatalf("expected action_count 1, got %d", p.ActionCount)
	}

	if err := b.QueuePlayerAction("missing", *a); err == nil {
		t.Fatalf("expected queue failure for missing player")
	}
}

func TestBoard_SnapshotCopies(t *testing.T) {
	p := &Player{ID: "p1", Name: "one", Health: 100, Actions: make(chan Action, PlayerRoundMoves)}
	c := &Creep{ID: "c1", Health: 100}
	u := &PowerUp{ID: "u1"}

	b := &Board{
		Name:     "test",
		Status:   BoardStatusRunning,
		SizeX:    20,
		SizeY:    20,
		Players:  []*Player{p},
		Creeps:   []*Creep{c},
		PowerUps: []*PowerUp{u},
		Events:   []*storage.Event{{Msg: "x"}},
	}

	info := b.SnapshotGameInfo()
	if info.AlivePlayers != 1 || info.AliveCreep != 1 || info.PowerUps != 1 {
		t.Fatalf("unexpected snapshot counts: %+v", info)
	}

	players := b.SnapshotPlayers()
	creeps := b.SnapshotCreeps()
	powerups := b.SnapshotPowerUps()
	events := b.SnapshotEvents()

	players[0] = nil
	creeps[0] = nil
	powerups[0] = nil
	events[0] = nil

	if len(b.SnapshotPlayers()) != 1 || b.SnapshotPlayers()[0] == nil {
		t.Fatalf("board players snapshot should be copy-backed")
	}
	if len(b.SnapshotCreeps()) != 1 || b.SnapshotCreeps()[0] == nil {
		t.Fatalf("board creeps snapshot should be copy-backed")
	}
	if len(b.SnapshotPowerUps()) != 1 || b.SnapshotPowerUps()[0] == nil {
		t.Fatalf("board powerups snapshot should be copy-backed")
	}
	if len(b.SnapshotEvents()) != 1 || b.SnapshotEvents()[0] == nil {
		t.Fatalf("board events snapshot should be copy-backed")
	}
}

func TestInitMetrics_Idempotent(t *testing.T) {
	InitMetrics()
	InitMetrics()
}

func TestNewBoard_RegistryRoundTrip(t *testing.T) {
	resetGameTestEnv(t)

	b := NewBoard("board-1", "name-1")
	CreateGame("board-1", b)

	got, ok := GetGame("board-1")
	if !ok {
		t.Fatalf("expected board in registry")
	}
	if got.ID != "board-1" || got.Name != "name-1" {
		t.Fatalf("unexpected board values: %+v", got)
	}
}
