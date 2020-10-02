package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sensepost/sconwar/game"
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

	go game.Games[id].Run()

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
// @Success 200 {object} GameResponse
// @Failure 400 {object} ErrorResponse
// @Router /games/get/:uuid [get]
func getGame(c *gin.Context) {

	p := &struct {
		UUID string `uri:"uuid" binding:"required,uuid"`
	}{}

	if err := c.ShouldBindUri(&p); err != nil {
		c.JSON(http.StatusBadRequest, &ErrorResponse{
			Message: `failed to read uri param`,
			Error:   err,
		})
	}

	c.JSON(http.StatusOK, &GameResponse{
		Game: game.Games[p.UUID],
	})
}
