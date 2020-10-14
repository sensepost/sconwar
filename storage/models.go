package storage

import (
	"time"

	"gorm.io/gorm"
)

// Player is the player model
type Player struct {
	gorm.Model

	Name   string
	UUID   string
	Scores []PlayerGameScore
}

// PlayerGameScore records the score for a player per game played
type PlayerGameScore struct {
	gorm.Model

	PlayerID      uint   `json:"-"`
	BoardID       string `json:"board_id"`
	Position      int    `json:"position"`
	Score         int    `json:"score"`
	DamageTaken   int    `json:"damage_taken"`
	DamageDealt   int    `json:"damage_dealt"`
	CreepKilled   int    `json:"killed_creep"`
	PlayersKilled int    `json:"killed_players"`
}

// Board is the board model
type Board struct {
	gorm.Model

	Name    string
	Status  int
	UUID    string
	Created time.Time
	Started time.Time
	Ended   time.Time
	Events  []*Event
}

// Event is a game event
type Event struct {
	gorm.Model

	BoardID uint `json:"board_id"`

	Date time.Time `json:"date_created"`

	SrcEntity   int    `json:"src_entity"`
	SrcEntityID string `json:"src_entity_id"`
	SrcPos      int    `json:"src_pos"`

	DstEntity   int    `json:"dst_entity"`
	DstEntityID string `json:"dst_entity_id"`
	DstPos      int    `json:"dst_pos"`

	Action int    `json:"action"`
	Msg    string `json:"msg"`
}
