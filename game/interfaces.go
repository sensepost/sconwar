package game

// hasPosition indicates a struct that has a position
type hasPosition interface {
	GetPosition() (int, int)
}

// hasInRange indicates a struct that can calculate range
type hasInRange interface {
	IsInRangeOf(hasPosition) bool
}

// takesDamage indicates a struct that has health
type takesDamage interface {
	TakeDamange(int)
}

// usesPowerUps indicates that the entity can use power ups
type usesPowerUps interface {
	GetPowerUp(PowerUp)
}
