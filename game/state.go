package game

import "sync"

// Games represent all of the games currently initialised
var Games map[string]*Board
var gamesMu sync.RWMutex

// Setup initialises the Games variable
func Setup() {
	gamesMu.Lock()
	defer gamesMu.Unlock()
	Games = make(map[string]*Board)
}

// CreateGame adds a new game to the in-memory registry.
func CreateGame(id string, board *Board) {
	gamesMu.Lock()
	defer gamesMu.Unlock()
	Games[id] = board
}

// GetGame gets a game by id.
func GetGame(id string) (*Board, bool) {
	gamesMu.RLock()
	defer gamesMu.RUnlock()
	board, ok := Games[id]
	return board, ok
}

// ListGames returns a snapshot of all active game pointers.
func ListGames() []*Board {
	gamesMu.RLock()
	defer gamesMu.RUnlock()

	games := make([]*Board, 0, len(Games))
	for _, g := range Games {
		games = append(games, g)
	}

	return games
}
