package api

import "github.com/sensepost/sconwar/game"

type NewPlayerResponse struct {
	Created bool
	UUID    string
}

type PlayerStatusResponse struct {
	Player *game.Player
}

type PlayerSurroundingResponse struct {
	Creep   []*game.Creep
	Players []*game.Player
}
