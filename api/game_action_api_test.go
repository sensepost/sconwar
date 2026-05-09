package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sensepost/sconwar/game"
	"github.com/sensepost/sconwar/storage"
)

func resetAPITestEnv(t *testing.T) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)
	if err := os.Setenv("API_TOKEN", "test-token"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}

	dbPath := filepath.Join(t.TempDir(), "db.sqlite")
	if err := storage.InitDbPath(dbPath); err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	game.Setup()

	return SetupRouter()
}

func mustCreateDBPlayer(t *testing.T, name string) *storage.Player {
	t.Helper()

	p := &storage.Player{
		Name: name,
		UUID: uuid.New().String(),
	}
	storage.Storage.Get().Create(p)
	return p
}

func performJSONRequest(t *testing.T, r http.Handler, method string, path string, body interface{}, token string) *httptest.ResponseRecorder {
	t.Helper()

	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal body: %v", err)
		}
	}

	req, err := http.NewRequest(method, path, bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("token", token)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestStartGame_IdempotentWhenRunning(t *testing.T) {
	r := resetAPITestEnv(t)

	id := uuid.New().String()
	b := game.NewBoard(id, "running-board")
	p := mustCreateDBPlayer(t, "joined")
	b.JoinPlayer(game.NewPlayer(p))
	b.SetStatus(game.BoardStatusRunning)
	game.CreateGame(id, b)

	w := performJSONRequest(t, r, http.MethodPut, "/api/game/start/"+id, nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestStopGame_IdempotentWhenFinished(t *testing.T) {
	r := resetAPITestEnv(t)

	id := uuid.New().String()
	b := game.NewBoard(id, "finished-board")
	b.SetStatus(game.BoardStatusFinished)
	game.CreateGame(id, b)

	w := performJSONRequest(t, r, http.MethodPut, "/api/game/stop/"+id, nil, "test-token")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestJoinGame_ForbiddenWhenNotNew(t *testing.T) {
	r := resetAPITestEnv(t)

	id := uuid.New().String()
	b := game.NewBoard(id, "closed-board")
	b.SetStatus(game.BoardStatusRunning)
	game.CreateGame(id, b)

	p := mustCreateDBPlayer(t, "p1")
	req := &JoinPlayerRequest{
		GameID:   id,
		PlayerID: p.UUID,
	}

	w := performJSONRequest(t, r, http.MethodPost, "/api/game/join", req, "")
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestJoinGame_ForbiddenWhenDuplicatePlayer(t *testing.T) {
	r := resetAPITestEnv(t)

	id := uuid.New().String()
	b := game.NewBoard(id, "new-board")
	game.CreateGame(id, b)

	p := mustCreateDBPlayer(t, "p1")
	b.JoinPlayer(game.NewPlayer(p))

	req := &JoinPlayerRequest{
		GameID:   id,
		PlayerID: p.UUID,
	}

	w := performJSONRequest(t, r, http.MethodPost, "/api/game/join", req, "")
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestMoveAction_BadRequestWhenPlayerNotInGame(t *testing.T) {
	r := resetAPITestEnv(t)

	id := uuid.New().String()
	b := game.NewBoard(id, "running-board")
	b.SetStatus(game.BoardStatusRunning)
	game.CreateGame(id, b)

	p := mustCreateDBPlayer(t, "db-only-player")

	req := &ActionMoveRequest{
		GamePlayerIDs: ActionGamePlayerRequest{
			GameID:   id,
			PlayerID: p.UUID,
		},
		X: 3,
		Y: 3,
	}

	w := performJSONRequest(t, r, http.MethodPost, "/api/action/move", req, "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
}
