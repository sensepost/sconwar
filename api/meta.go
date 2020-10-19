package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sensepost/sconwar/game"
)

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
