package game

// Game defaults
const (
	BoardX                = 20
	BoardY                = 20
	CreepCount            = 15  // how many initial creep to spawn
	FogOfWarPercent       = 100 // todo: fix reduce prod; how far a player can see other entities
	AttackRange           = 2   // distance considered in range for an attack
	PlayerRoundMoves      = 2   // how many rounds a player gets to move
	CreepRoundMoves       = 3   // how many rounds a creep gets to move
	MaxRoundSeconds       = 30  // how long a player can take to issue actions
	MaxPlayerMoveDistance = 2   // how far a single move for a player can be
	MaxCreepMoveDistance  = 1   // how far a single move for a creep can be
	PowerUpMax            = 5   // the max number of powerups available on the board
	PowerUpChance         = 50  // the chance of a powerup spawning
	PowerUpHealthBonus    = 50  // how much extra health the health powerup gives
)
