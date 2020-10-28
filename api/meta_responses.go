package api

import (
	"github.com/sensepost/sconwar/game"
)

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

// MetaTotalScoresResponse is a response containing game types
type MetaTotalScoresResponse struct {
	Players []*PlayerTotalScore `json:"players"`
}

// PlayerTotalScore contains player score totals
type PlayerTotalScore struct {
	Name             string `json:"name"`
	AveragePosition  int    `json:"average_position"`
	TotalScore       int    `json:"total_score"`
	TotalDamageTaken int    `json:"total_damage_taken"`
	TotalDamageDealt int    `json:"total_damage_dealt"`
	TotalCreepKills  int    `json:"total_creep_kills"`
	TotalPlayerKills int    `json:"total_player_kills"`
}

// PlayerLeaderBoardResponse contains the leader board
type PlayerLeaderBoardResponse struct {
	Scores []*PlayerLeaderboardScore `json:"scores"`
}

// PlayerLeaderboardScore is a player score representation
type PlayerLeaderboardScore struct {
	Name        string `json:"name"`
	GameID      string `json:"game_id"`
	GameName    string `json:"game_name"`
	Score       int    `json:"score"`
	Position    int    `json:"position"`
	DamageTaken int    `json:"damage_taken"`
	DamageDealt int    `json:"damage_dealt"`
	CreepKills  int    `json:"creep_kills"`
	PlayerKills int    `json:"player_kills"`
}
