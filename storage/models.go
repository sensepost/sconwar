package storage

import "gorm.io/gorm"

// Player contains information about a Player
type Player struct {
	gorm.Model

	Name string
	UUID string
}
