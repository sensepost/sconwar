package game

import "github.com/prometheus/client_golang/prometheus"

var PlayerApiActions = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "sconwar_player_api_actions_total",
		Help: "number of different actions performed using the player api",
	},
	[]string{"action"},
)


var attackActionExecuted = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "sconwar_attack_executed_total",
		Help: "count of attacks executed",
	},
	[]string{},
)

var playerCreated = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "sconwar_player_created_total",
		Help: "count of all players created",
	},
	[]string{},
)

var distanceMovedByEntity = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "sconwar_distance_moved_total",
		Help: "total distance moved by entities",
	},
	[]string{"entity"},
)

var playerActionsQueued = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "sconwar_player_actions_queued_total",
		Help: "count of all actions queued",
	},
	[]string{},
)

var damageTaken = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "sconwar_damage_taken_total",
		Help: "total amount of damage taken by entities",
	},
	[]string{"entity"},
)

var powerupsCollected = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "sconwar_powerups_collected_total",
		Help: "total number of powerups collected",
	},
	[]string{},
)

var powerupsUsed = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "sconwar_powerups_used_total",
		Help: "total number of powerups used",
	},
	[]string{"poweruptype"},
)

var scoreAwarded = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "sconwar_score_awarded_total",
		Help: "total score awarded to players",
	},
	[]string{},
)

var creepsKilled = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "sconwar_creeps_killed_total",
		Help: "total number of creeps killed",
	},
	[]string{},
)

var playersKilled = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "sconwar_players_killed_total",
		Help: "total number of players killed",
	},
	[]string{},
)

var gameState = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "sconwar_game_state_total",
		Help: "total number of game states executed (new/finished)",
	},
	[]string{"state"},
)

//InitMetrics sets up the prometheus metrics for exposing over GIN
func InitMetrics() {
	prometheus.Register(attackActionExecuted)
	prometheus.Register(playerCreated)
	prometheus.Register(playerActionsQueued)
	prometheus.Register(damageTaken)
	prometheus.Register(powerupsCollected)
	prometheus.Register(powerupsUsed)
	prometheus.Register(scoreAwarded)
	prometheus.Register(creepsKilled)
	prometheus.Register(playersKilled)

	prometheus.Register(PlayerApiActions)
}
