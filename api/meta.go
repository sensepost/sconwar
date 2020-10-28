package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sensepost/sconwar/game"
	"github.com/sensepost/sconwar/storage"
)

// GetTotalScores godoc
// @Summary Get the total scores for everyone
// @Description Get's total scores for everyone for every game that has been played
// @Tags Meta
// @Accept json
// @Produce json
// @Success 200 {object} MetaTotalScoresResponse
// @Router /meta/scores/ [get]
func getTotalScores(c *gin.Context) {

	res := &MetaTotalScoresResponse{}

	var players []storage.Player
	storage.Storage.Get().Preload("Scores").Find(&players)

	for _, p := range players {

		score := &PlayerTotalScore{
			Name: p.Name,
		}

		for _, ps := range p.Scores {
			score.TotalScore += ps.Score
			score.TotalDamageTaken += ps.DamageTaken
			score.TotalDamageDealt += ps.DamageDealt
			score.TotalCreepKills += ps.CreepKilled
			score.TotalPlayerKills += ps.PlayersKilled
			score.AveragePosition += ps.Position
		}

		score.AveragePosition = score.AveragePosition / len(p.Scores)

		res.Players = append(res.Players, score)
	}

	c.JSON(http.StatusOK, res)
}

// GetLeaderBoard godoc
// @Summary Get the leader board
// @Description Get's the leader board for this game server
// @Tags Meta
// @Accept json
// @Produce json
// @Success 200 {object} PlayerLeaderBoardResponse
// @Router /meta/leaderboard/ [get]
func getLeaderBoard(c *gin.Context) {

	res := &PlayerLeaderBoardResponse{}

	var scores []storage.PlayerGameScore
	storage.Storage.Get().Order("score desc").Limit(20).Find(&scores)

	for _, s := range scores {
		var player storage.Player
		storage.Storage.Get().Find(&player, s.PlayerID)

		var board storage.Board
		storage.Storage.Get().Where("uuid = ?", s.BoardID).First(&board)

		score := &PlayerLeaderboardScore{
			Name:        player.Name,
			GameID:      s.BoardID,
			GameName:    board.Name,
			Score:       s.Score,
			Position:    s.Position,
			DamageTaken: s.DamageTaken,
			DamageDealt: s.DamageDealt,
			CreepKills:  s.CreepKilled,
			PlayerKills: s.PlayersKilled,
		}

		res.Scores = append(res.Scores, score)
	}

	c.JSON(http.StatusOK, res)
}

// GetTypes godoc
// @Summary Get game type information
// @Description Get's information about types & enumerations in the game
// @Tags Meta
// @Accept json
// @Produce json
// @Success 200 {object} MetaTypesResponse
// @Router /meta/types/ [get]
func getTypes(c *gin.Context) {

	c.JSON(http.StatusOK, &MetaTypesResponse{
		PlayerActions: &PlayerActions{
			Move:    game.Move,
			Attack:  game.Attack,
			Pickup:  game.Pickup,
			Nothing: game.Nothing,
		},
		BoardStatuses: &BoardStatuses{
			New:      game.BoardStatusNew,
			Running:  game.BoardStatusRunning,
			Finished: game.BoardStatusFinished,
		},
		GameEntities: &GameEntities{
			PlayerEntity:  game.PlayerEntity,
			CreepEntity:   game.CreepEntity,
			PowerupEntity: game.PowerupEntity,
		},
		PowerupTypes: &PowerupTypes{
			Health:       game.Health,
			Teleport:     game.Teleport,
			DoubleDamage: game.DoubleDamage,
		},
	})
}
