package api

import "github.com/sensepost/sconwar/game"

// NewPlayerResponse is a response for a newly registered player
type NewPlayerResponse struct {
	Created bool
	UUID    string
}

// PlayerStatusResponse is a player status response
type PlayerStatusResponse struct {
	Player *game.Player
}

// PlayerSurroundingResponse is a surroundings response
type PlayerSurroundingResponse struct {
	Creep   []*game.Creep
	Players []*game.Player
}
