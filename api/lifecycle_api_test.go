package api

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sensepost/sconwar/game"
)

func decodeBodyJSON(t *testing.T, body []byte, target interface{}) {
	t.Helper()
	if err := json.Unmarshal(body, target); err != nil {
		t.Fatalf("failed to decode json: %v body=%s", err, string(body))
	}
}

func waitForGameFinished(t *testing.T, r http.Handler, gameID string, timeout time.Duration) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		w := performJSONRequest(t, r, http.MethodGet, "/api/game/info/"+gameID, nil, "")
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 on game info, got %d body=%s", w.Code, w.Body.String())
		}

		var info GameInfoResponse
		decodeBodyJSON(t, w.Body.Bytes(), &info)
		if info.Status == game.BoardStatusFinished {
			return
		}

		time.Sleep(25 * time.Millisecond)
	}

	t.Fatalf("game %s did not finish within %s", gameID, timeout)
}

func TestAPIIntegration_LifecycleAndPersistence(t *testing.T) {
	r := resetAPITestEnv(t)

	// register player
	registerReq := &RegisterPlayerRequest{Name: "integration-player"}
	w := performJSONRequest(t, r, http.MethodPost, "/api/player/register", registerReq, "test-token")
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 register, got %d body=%s", w.Code, w.Body.String())
	}

	var reg NewPlayerResponse
	decodeBodyJSON(t, w.Body.Bytes(), &reg)
	if reg.UUID == "" {
		t.Fatalf("expected player uuid in register response")
	}

	// create game
	newGameReq := &NewGameRequest{Name: "integration-game"}
	w = performJSONRequest(t, r, http.MethodPost, "/api/game/new", newGameReq, "")
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 new game, got %d body=%s", w.Code, w.Body.String())
	}

	var ng NewGameResponse
	decodeBodyJSON(t, w.Body.Bytes(), &ng)
	if ng.UUID == "" {
		t.Fatalf("expected game uuid in new game response")
	}

	// make deterministic/fast: no creep/powerups so game can end after queued player turn.
	board, ok := game.GetGame(ng.UUID)
	if !ok {
		t.Fatalf("expected game to exist in registry")
	}
	board.Creeps = nil
	board.PowerUps = nil

	// join game
	joinReq := &JoinPlayerRequest{
		GameID:   ng.UUID,
		PlayerID: reg.UUID,
	}
	w = performJSONRequest(t, r, http.MethodPost, "/api/game/join", joinReq, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 join game, got %d body=%s", w.Code, w.Body.String())
	}

	// create non-member db player and ensure action rejected for membership.
	other := mustCreateDBPlayer(t, "db-not-in-game")
	badMove := &ActionMoveRequest{
		GamePlayerIDs: ActionGamePlayerRequest{
			GameID:   ng.UUID,
			PlayerID: other.UUID,
		},
		X: 1,
		Y: 1,
	}
	w = performJSONRequest(t, r, http.MethodPost, "/api/action/move", badMove, "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for non-member action, got %d body=%s", w.Code, w.Body.String())
	}

	// start game
	w = performJSONRequest(t, r, http.MethodPut, "/api/game/start/"+ng.UUID, nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 start, got %d body=%s", w.Code, w.Body.String())
	}

	// queue two actions to satisfy the player round immediately.
	move1 := &ActionMoveRequest{
		GamePlayerIDs: ActionGamePlayerRequest{
			GameID:   ng.UUID,
			PlayerID: reg.UUID,
		},
		X: 2,
		Y: 2,
	}
	move2 := &ActionMoveRequest{
		GamePlayerIDs: ActionGamePlayerRequest{
			GameID:   ng.UUID,
			PlayerID: reg.UUID,
		},
		X: 3,
		Y: 3,
	}
	w = performJSONRequest(t, r, http.MethodPost, "/api/action/move", move1, "")
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 action move1, got %d body=%s", w.Code, w.Body.String())
	}
	w = performJSONRequest(t, r, http.MethodPost, "/api/action/move", move2, "")
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 action move2, got %d body=%s", w.Code, w.Body.String())
	}

	waitForGameFinished(t, r, ng.UUID, 3*time.Second)

	// event stream should have content
	w = performJSONRequest(t, r, http.MethodGet, "/api/game/events/"+ng.UUID, nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 game events, got %d body=%s", w.Code, w.Body.String())
	}
	var events GameEventsResponse
	decodeBodyJSON(t, w.Body.Bytes(), &events)
	if len(events.Events) == 0 {
		t.Fatalf("expected at least one event in game")
	}

	// scores endpoint should include player
	w = performJSONRequest(t, r, http.MethodGet, "/api/game/scores/"+ng.UUID, nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 game scores, got %d body=%s", w.Code, w.Body.String())
	}
	var scores GameScoresResponse
	decodeBodyJSON(t, w.Body.Bytes(), &scores)
	if len(scores.Scores) == 0 {
		t.Fatalf("expected at least one score row")
	}

	// meta totals should include our player by name.
	w = performJSONRequest(t, r, http.MethodGet, "/api/meta/scores", nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 meta scores, got %d body=%s", w.Code, w.Body.String())
	}
	var totals MetaTotalScoresResponse
	decodeBodyJSON(t, w.Body.Bytes(), &totals)
	found := false
	for _, p := range totals.Players {
		if p.Name == "integration-player" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected integration-player in meta total scores")
	}
}

func TestAPIIntegration_SecuredEndpointRejectsInvalidToken(t *testing.T) {
	r := resetAPITestEnv(t)

	id := uuid.New().String()
	b := game.NewBoard(id, "token-check")
	game.CreateGame(id, b)

	w := performJSONRequest(t, r, http.MethodPut, "/api/game/stop/"+id, nil, "wrong-token")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for invalid token, got %d body=%s", w.Code, w.Body.String())
	}
}
