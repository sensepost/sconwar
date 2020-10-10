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
	UUID    string
	Created time.Time
	Started time.Time
	Events  []*Event
}

// Event is a game event
type Event struct {
	gorm.Model

	BoardID uint

	Date time.Time

	SrcEntity   int
	SrcEntityID string
	SrcPos      int

	DstEntity   int
	DstEntityID string
	DstPos      int

	Action int
	Msg    string
}
