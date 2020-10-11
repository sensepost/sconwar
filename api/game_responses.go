package api

import (
	"time"

	"github.com/sensepost/sconwar/game"
	"github.com/sensepost/sconwar/storage"
)

// NewGameResponse is a new game response
type NewGameResponse struct {
	Created bool   `json:"created"`
	UUID    string `json:"uuid"`
}

// AllGamesResponse is a response with all games
type AllGamesResponse struct {
	Games []string `json:"games"`
}

// GameDetailResponse is a response with a single game
type GameDetailResponse struct {
	Game *game.Board `json:"game"`
}

// GameInfoResponse is a response summary for a game
type GameInfoResponse struct {
	Name    string    `json:"name"`
	SizeX   int       `json:"size_x"`
	SizeY   int       `json:"size_y"`
	Fow     float64   `json:"fow"`
	Created time.Time `json:"created"`
	Started time.Time `json:"started"`
}

// EventsResponse is a response with a games' events
type EventsResponse struct {
	Events []*storage.Event `json:"events"`
}
