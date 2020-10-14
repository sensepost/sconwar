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
	Name          string               `json:"name"`
	Status        game.BoardStatus     `json:"status"`
	SizeX         int                  `json:"size_x"`
	SizeY         int                  `json:"size_y"`
	Fow           float64              `json:"fow"`
	CurrentPlayer string               `json:"current_player"`
	Created       time.Time            `json:"created"`
	Started       time.Time            `json:"started"`
	GameOptions   GameOptionsResponse  `json:"game_options"`
	GameEntities  GameEntitiesResponse `json:"game_entities"`
}

// GameOptionsResponse are the options in the GameInfoResponse
type GameOptionsResponse struct {
	FogOfWarPercent       int `json:"fow_percent"`
	AttackRange           int `json:"player_attack_range"`
	PlayerRoundMoves      int `json:"player_round_moves"`
	MaxRoundSeconds       int `json:"player_round_seconds"`
	MaxPlayerMoveDistance int `json:"player_max_move_range"`
	PowerUpMax            int `json:"powerup_max"`
}

// GameEntitiesResponse contains the number of enties in a game
type GameEntitiesResponse struct {
	AliveCreep   int `json:"alive_creep"`
	AlivePlayers int `json:"alive_players"`
	PowerUps     int `json:"powerups"`
}

// GameEventsResponse is a response with a games' events
type GameEventsResponse struct {
	Events []*storage.Event `json:"events"`
}

// GameScoresResponse is a response with a games' events
type GameScoresResponse struct {
	Scores []*PlayerScore `json:"scores"`
}

// PlayerScore contains a players score
type PlayerScore struct {
	Name          string `json:"name"`
	Score         int    `json:"score"`
	DamageDealt   int    `json:"damage_dealt"`
	DamageTaken   int    `json:"damage_taken"`
	CreepKilled   int    `json:"killed_creep"`
	PlayersKilled int    `json:"killed_players"`
}
