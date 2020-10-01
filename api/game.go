package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sensepost/sconwar/game"
)

func newGame(c *gin.Context) {

	id := uuid.New().String()
	game.Games[id] = game.NewBoard(id)

	go game.Games[id].Run()

	c.JSON(http.StatusOK, gin.H{
		"created": true,
		"uuid":    id,
	})
}

func allGames(c *gin.Context) {
	c.JSON(http.StatusOK, game.Games)
}
