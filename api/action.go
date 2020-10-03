package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sensepost/sconwar/game"
)

// MoveAction godoc
// @Summary Move a player in a gme
// @Description Move's a player in a game to a new co-ordinate
// @Tags Action
// @Accept json
// @Produce json
// @Param data body	ActionMoveRequest true "ActionMoveRequest Request"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Router /action/move [post]
func moveAction(c *gin.Context) {

	params := &ActionMoveRequest{}

	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{
			Message: `failed to read post params`,
			Error:   err.Error(),
		})
		return
	}

	if err := params.Validation(); err != nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{
			Message: `failed to validate request`,
			Error:   err.Error(),
		})
		return
	}

	var player *game.Player
	for _, p := range game.Games[params.GamePlayerIDs.GameID].Players {
		if p.ID == params.GamePlayerIDs.PlayerID {
			player = p
		}
	}

	action := game.NewAction(*player, game.Move)
	action.SetXY(params.X, params.Y)

	if err := player.AddAction(*action); err != nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{
			Message: `failed to buffer new action`,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &StatusResponse{
		Success: true,
	})
}
