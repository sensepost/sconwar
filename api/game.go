package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sensepost/sconwar/game"
	"github.com/sensepost/sconwar/storage"
)

// NewGame godoc
// @Summary Register a new game
// @Description Registers the start of a new game, returning the game id
// @Tags Games
// @Accept json
// @Produce json
// @Success 200 {object} NewGameResponse
// @Router /games/new [get]
func newGame(c *gin.Context) {

	id := uuid.New().String()
	game.Games[id] = game.NewBoard(id)

	c.JSON(http.StatusOK, &NewGameResponse{
		Created: true,
		UUID:    id,
	})
}

// AllGames godoc
// @Summary List all games
// @Description Get's a list of all of the games that are running
// @Tags Games
// @Accept json
// @Produce json
// @Success 200 {object} AllGamesResponse
// @Router /games/ [get]
func allGames(c *gin.Context) {

	g := &AllGamesResponse{}
	for _, d := range game.Games {
		g.Games = append(g.Games, d.ID)
	}

	c.JSON(http.StatusOK, g)
}

// GetGame godoc
// @Summary Get game details
// @Description Get's the details for a game defined by UUID
// @Tags Games
// @Accept json
// @Produce json
// @Param game_id path string true "game uuid"
// @Success 200 {object} GameResponse
// @Failure 400 {object} ErrorResponse
// @Router /games/get/{game_id} [get]
func getGame(c *gin.Context) {

	params := &GetGameDetailRequest{}

	if err := c.BindUri(&params); err != nil {
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

	c.JSON(http.StatusOK, &GameResponse{
		Game: game.Games[params.GameID],
	})
}

// StartGame godoc
// @Summary Start a game
// @Description Starts a game
// @Tags Games
// @Accept json
// @Produce json
// @Param game_id path string true "game uuid"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Router /games/start/{game_id} [put]
func startGame(c *gin.Context) {

	params := &GetGameDetailRequest{}

	if err := c.BindUri(&params); err != nil {
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

	// start the game
	go game.Games[params.GameID].Run()

	c.JSON(http.StatusOK, &StatusResponse{
		Success: true,
	})
}

// JoinGame godoc
// @Summary Join a game
// @Description Joins a player to an existing game
// @Tags Games
// @Accept json
// @Produce json
// @Param data body	JoinPlayerRequest true "Join Request"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Router /games/join [post]
func joinGame(c *gin.Context) {

	// todo: check that player is not already in the game
	// todo: check that the game is not already running

	params := &JoinPlayerRequest{}

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
	storage.Storage.Get().Where("UUID = ?", params.PlayerID).First(&player)

	gamePlayer := game.NewPlayer(&player)
	game.Games[params.GameID].JoinPlayer(gamePlayer)

	c.JSON(http.StatusOK, &StatusResponse{
		Success: true,
	})
}
