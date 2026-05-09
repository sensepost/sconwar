package api

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sensepost/sconwar/game"
	"github.com/sensepost/sconwar/storage"
)

// GetPlayer godoc
// @Summary Get player information
// @Description Get's information about a player by UUID
// @Tags Player
// @Accept json
// @Produce json
// @Param data body	PlayerRequest true "PlayerRequest Request"
// @Success 200 {object} storage.Player
// @Failure 400 {object} ErrorResponse
// @Router /player/ [post]
func getPlayer(c *gin.Context) {

	params := &PlayerRequest{}

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

	var player storage.Player
	storage.Storage.Get().Where("uuid = ?", params.PlayerID).Preload("Scores").First(&player)

	game.PlayerApiActions.With(prometheus.Labels{"action": "info"}).Inc()
	c.JSON(http.StatusOK, player)
}

// RegisterPlayer godoc
// @Summary Register a new player
// @Description Registers a new player that can join games. The returned UUID is the secret too!
// @Tags Player
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body	RegisterPlayerRequest true "RegisterPlayerRequest Request"
// @Success 201 {object} NewPlayerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
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

	game.PlayerApiActions.With(prometheus.Labels{"action": "register"}).Inc()
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

	board, _ := game.GetGame(params.GameID)
	status, _ := board.FindPlayer(params.PlayerID)

	game.PlayerApiActions.With(prometheus.Labels{"action": "status"}).Inc()
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

	distances := &PlayerSurroundingResponse{}
	board, _ := game.GetGame(params.GameID)
	player, _ := board.FindPlayer(params.PlayerID)

	// if the board is complete (aka, game done) we disable fow calculation
	if board.StatusValue() == game.BoardStatusFinished {
		distances.FOWEnabled = false

		distances.PowerUps = board.SnapshotPowerUps()
		distances.Creep = board.SnapshotCreeps()
		distances.Players = board.SnapshotPlayers()
	} else {

		distances.FOWEnabled = true

		info := board.SnapshotGameInfo()
		for _, u := range board.SnapshotPowerUps() {
			if player.DistanceFrom(u) <= info.FOWDistance {
				distances.PowerUps = append(distances.PowerUps, u)
			}
		}

		for _, c := range board.SnapshotCreeps() {
			if player.DistanceFrom(c) <= info.FOWDistance && c.Health > 0 {
				distances.Creep = append(distances.Creep, c)
			}
		}

		for _, p := range board.SnapshotPlayers() {
			if player.ID == p.ID {
				continue
			}

			if player.DistanceFrom(p) <= info.FOWDistance && p.Health > 0 {
				distances.Players = append(distances.Players, p)
			}
		}

	}

	game.PlayerApiActions.With(prometheus.Labels{"action": "surroundings"}).Inc()
	c.JSON(http.StatusOK, &distances)
}
