package game

// Entity is a game entity
type Entity int

// Game entities that can be logged
const (
	PlayerEntity Entity = iota
	CreepEntity
	PowerupEntity
)

// the model for an event lives in the storage package
