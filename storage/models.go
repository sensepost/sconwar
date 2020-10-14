package storage

import (
	"time"

	"gorm.io/gorm"
)

// Player is the player model
type Player struct {
	gorm.Model

	Name string
	UUID string
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
