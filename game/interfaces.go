package game

type hasPosition interface {
	GetPosition() (int, int)
}

type hasInRange interface {
	IsInRangeOf(hasPosition) bool
}

type takesDamage interface {
	TakeDamange(int)
}
