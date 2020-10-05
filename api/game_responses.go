package api

import "github.com/sensepost/sconwar/game"

// NewGameResponse is a new game response
type NewGameResponse struct {
	Created bool   `json:"created"`
	UUID    string `json:"uuid"`
}

// AllGamesResponse is a response with all games
type AllGamesResponse struct {
	Games []string `json:"games"`
}

// GameResponse is a response with a single game
type GameResponse struct {
	Game *game.Board `json:"game"`
}
