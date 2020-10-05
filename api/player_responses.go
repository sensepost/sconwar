package api

import "github.com/sensepost/sconwar/game"

// NewPlayerResponse is a response for a newly registered player
type NewPlayerResponse struct {
	Created bool   `json:"created"`
	UUID    string `json:"uuid"`
}

// PlayerStatusResponse is a player status response
type PlayerStatusResponse struct {
	Player *game.Player `json:"player"`
}

// PlayerSurroundingResponse is a surroundings response
type PlayerSurroundingResponse struct {
	Creep    []*game.Creep   `json:"creep"`
	Players  []*game.Player  `json:"players"`
	PowerUps []*game.PowerUp `json:"powerups"`
}
