package api

import "github.com/sensepost/sconwar/game"

// MetaTypesResponse is a response containing game types
type MetaTypesResponse struct {
	PlayerActions *PlayerActions `json:"player_actions"`
	BoardStatuses *BoardStatuses `json:"board_statuses"`
	GameEntities  *GameEntities  `json:"game_entities"`
	PowerupTypes  *PowerupTypes  `json:"powerup_types"`
}

// PlayerActions are the player actions
type PlayerActions struct {
	Move    game.ActionType `json:"move"`
	Attack  game.ActionType `json:"attack"`
	Pickup  game.ActionType `json:"pickup"`
	Nothing game.ActionType `json:"nothing"`
}

// BoardStatuses are the game board statuses
type BoardStatuses struct {
	New      game.BoardStatus `json:"new"`
	Running  game.BoardStatus `json:"running"`
	Finished game.BoardStatus `json:"finished"`
}

// GameEntities are the game entities
type GameEntities struct {
	PlayerEntity  game.Entity `json:"player"`
	CreepEntity   game.Entity `json:"creep"`
	PowerupEntity game.Entity `json:"powerup"`
}

// PowerupTypes are the powerups
type PowerupTypes struct {
	Health       game.PowerUpType `json:"health"`
	Teleport     game.PowerUpType `json:"teleport"`
	DoubleDamage game.PowerUpType `json:"double_damage"`
}
