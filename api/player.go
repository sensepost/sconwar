package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sensepost/sconwar/storage"
	"gorm.io/gorm"
)

// RegisterPlayer godoc
// @Summary Register a new player
// @Description Registers a new player that can join games. The returned UUID is the secret too!
// @Tags Player
// @Accept json
// @Produce json
// @Param data body	RegisterPlayer true "RegsiterPlayer Request"
// @Success 200 {object} NewPlayerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /player/register [post]
func registerPlayer(c *gin.Context) {

	params := &RegisterPlayer{}

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

	// does the player name already exist?
	var player storage.Player
	res := storage.Storage.Get().Where("name = ?", params.Name).First(&player)
	if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusForbidden, &ErrorResponse{
			Message: `player name already registered`,
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

	c.JSON(http.StatusOK, &NewPlayerResponse{
		Created: true,
		UUID:    id,
	})
}
