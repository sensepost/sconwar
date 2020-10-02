package api

import "github.com/sensepost/sconwar/game"

type ErrorResponse struct {
	Message string
	Error   string
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

type StatusResponse struct {
	Success bool
}
