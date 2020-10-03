package api

import "github.com/sensepost/sconwar/game"

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
