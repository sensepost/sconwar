package api

import "github.com/sensepost/sconwar/game"

type ErrorResponse struct {
	Message string
	Error   error
}

type NewGameResponse struct {
	Created bool
	UUID    string
}

type AllGamesResponse struct {
	Games []string
}

type GameResponse struct {
	Game *game.Board
}
