package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sensepost/sconwar/game"
	"github.com/sensepost/sconwar/storage"
)

// RegisterPlayer godoc
// @Summary Register a new player
// @Description Registers a new player that can join games. The returned UUID is the secret too!
// @Tags Player
// @Accept json
// @Produce json
// @Param data body	RegisterPlayerRequest true "RegisterPlayerRequest Request"
// @Success 201 {object} NewPlayerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /player/register [post]
func registerPlayer(c *gin.Context) {

	params := &RegisterPlayerRequest{}

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

	// create the player
	id := uuid.New().String()
	newPlayer := &storage.Player{
		Name: params.Name,
		UUID: id,
	}

	storage.Storage.Get().Create(newPlayer)

	c.JSON(http.StatusCreated, &NewPlayerResponse{
		Created: true,
		UUID:    id,
	})
}

// PlayerStatus godoc
// @Summary Get Player status in a game
// @Description Get's the player status in a game
// @Tags Player
// @Accept json
// @Produce json
// @Param data body	PlayerGameRequest true "PlayerGameRequest Request"
// @Success 200 {object} PlayerStatusResponse
// @Failure 400 {object} ErrorResponse
// @Router /player/status [post]
func playerStatus(c *gin.Context) {

	params := &PlayerGameRequest{}

	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{
			Message: `failed to read uri param`,
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

	var status *game.Player
	for _, p := range game.Games[params.GameID].Players {
		if p.ID == params.PlayerID {
			status = p
		}
	}

	c.JSON(http.StatusOK, &PlayerStatusResponse{
		Player: status,
	})
}

// PlayerSurrounding godoc
// @Summary Get a player's surroundings
// @Description Get's the surroundings of a player in a game, outside of fog of war
// @Tags Player
// @Accept json
// @Produce json
// @Param data body	PlayerGameRequest true "PlayerGameRequest Request"
// @Success 200 {object} PlayerSurroundingResponse
// @Failure 400 {object} ErrorResponse
// @Router /player/surroundings [post]
func playerSurrounding(c *gin.Context) {

	params := &PlayerGameRequest{}

	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{
			Message: `failed to read uri param`,
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
	for _, p := range game.Games[params.GameID].Players {
		if p.ID == params.PlayerID {
			player = p
		}
	}

	distances := &PlayerSurroundingResponse{}
	board := game.Games[params.GameID]

	for _, u := range board.PowerUps {
		if player.DistanceFrom(u) <= board.FOWDistance {
			distances.PowerUps = append(distances.PowerUps, u)
		}
	}

	for _, c := range board.Creeps {
		if player.DistanceFrom(c) <= board.FOWDistance && c.Health > 0 {
			distances.Creep = append(distances.Creep, c)
		}
	}

	for _, p := range board.Players {
		if player.ID == p.ID {
			continue
		}

		if player.DistanceFrom(p) <= board.FOWDistance && p.Health > 0 {
			distances.Players = append(distances.Players, p)
		}
	}

	c.JSON(http.StatusOK, &distances)
}
