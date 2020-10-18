package game

// Game defaults
const (
	BoardX                = 20
	BoardY                = 20
	CreepCount            = 15  // how many initial creep to spawn
	FogOfWarPercent       = 30  // how far a player can see other entities
	AttackRange           = 2   // distance considered in range for an attack
	PlayerRoundMoves      = 2   // how many rounds a player gets to move
	CreepRoundMoves       = 2   // how many rounds a creep gets to move
	MaxRoundSeconds       = 30  // how long a player can take to issue actions
	MaxPlayerMoveDistance = 2   // how far a single move for a player can be
	MaxCreepMoveDistance  = 1   // how far a single move for a creep can be
	MaxDamage             = 20  // how much damage can one take, at most
	PowerUpMax            = 5   // the max number of powerups available on the board
	PowerUpChance         = 50  // the chance of a powerup spawning
	PowerUpHealthBonus    = 50  // how much extra health the health powerup gives
	PowerUpHealthBonusMax = 135 // the max amount of health a player can have

	// scoring
	CreepKilledScore  = 100 // how many points to add for killing a creep
	PlayerKilledScore = 200 // how many points to add for killing a player
	PickedUpPowerup   = 50  // how many points to add for picking up a powerup
)
