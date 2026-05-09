package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sensepost/sconwar/game"
)

// MoveAction godoc
// @Summary Move a player in a game
// @Description Move's a player in a game to a new co-ordinate
// @Tags Action
// @Accept json
// @Produce json
// @Param data body	ActionMoveRequest true "ActionMoveRequest Request"
// @Success 201 {object} StatusResponse
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

	board, _ := game.GetGame(params.GamePlayerIDs.GameID)

	action := game.NewAction(game.Move)
	action.SetXY(params.X, params.Y)

	if err := board.QueuePlayerAction(params.GamePlayerIDs.PlayerID, *action); err != nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{
			Message: `failed to buffer new action`,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &StatusResponse{
		Success: true,
	})
}

// AttackAction godoc
// @Summary Attack an entity
// @Description Attack's an entity at x, y, assuming it's in range
// @Tags Action
// @Accept json
// @Produce json
// @Param data body	ActionAttackRequest true "ActionAttackRequest Request"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Router /action/attack [post]
func attackAction(c *gin.Context) {

	params := &ActionAttackRequest{}

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

	board, _ := game.GetGame(params.GamePlayerIDs.GameID)

	action := game.NewAction(game.Attack)
	action.SetXY(params.X, params.Y)

	if err := board.QueuePlayerAction(params.GamePlayerIDs.PlayerID, *action); err != nil {
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

// PickupAction godoc
// @Summary Pick up an item
// @Description Pick's an item up and places it in the players inventory
// @Tags Action
// @Accept json
// @Produce json
// @Param data body	ActionPickupRequest true "ActionPickupRequest Request"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Router /action/pickup [post]
func pickupAction(c *gin.Context) {

	params := &ActionPickupRequest{}

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

	board, _ := game.GetGame(params.GamePlayerIDs.GameID)

	action := game.NewAction(game.Pickup)
	action.SetXY(params.X, params.Y)

	if err := board.QueuePlayerAction(params.GamePlayerIDs.PlayerID, *action); err != nil {
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

// UseAction godoc
// @Summary Use a powerup
// @Description Uses's a powerup and activates it buff
// @Tags Action
// @Accept json
// @Produce json
// @Param data body	ActionUseRequest true "ActionUseRequest Request"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Router /action/use [post]
func useAction(c *gin.Context) {

	params := &ActionUseRequest{}

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

	board, _ := game.GetGame(params.GamePlayerIDs.GameID)
	player, ok := board.FindPlayer(params.GamePlayerIDs.PlayerID)
	if !ok {
		c.JSON(http.StatusBadRequest, &ErrorResponse{
			Message: `failed to validate request`,
			Error:   `this player is not part of this game`,
		})
		return
	}

	player.UsePowerUp(params.PowerUpID)

	c.JSON(http.StatusOK, &StatusResponse{
		Success: true,
	})
}
