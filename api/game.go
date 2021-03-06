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
// @Tags Game
// @Accept json
// @Produce json
// @Param data body	NewGameRequest true "NewGameRequest"
// @Success 201 {object} NewGameResponse
// @Failure 400 {object} ErrorResponse
// @Router /game/new [post]
func newGame(c *gin.Context) {

	params := &NewGameRequest{}

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

	id := uuid.New().String()
	game.Games[id] = game.NewBoard(id, params.Name)

	c.JSON(http.StatusCreated, &NewGameResponse{
		Created: true,
		UUID:    id,
	})
}

// AllGames godoc
// @Summary List all games
// @Description Get's a list of all of the games that are running
// @Tags Game
// @Accept json
// @Produce json
// @Success 200 {object} AllGamesResponse
// @Router /game/ [get]
func allGames(c *gin.Context) {
	g := &AllGamesResponse{}
	for _, d := range game.Games {
		g.Games = append(g.Games, &AllGamesGame{
			ID:     d.ID,
			Name:   d.Name,
			Status: int(d.Status),
		})
	}

	c.JSON(http.StatusOK, g)
}

// GetGameDetail godoc
// @Summary Get game details
// @Description Get's the details for a game defined by UUID
// @Tags Game
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param game_id path string true "game uuid"
// @Success 200 {object} GameDetailResponse
// @Failure 400 {object} ErrorResponse
// @Router /game/detail/{game_id} [get]
func getGameDetail(c *gin.Context) {

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

	c.JSON(http.StatusOK, &GameDetailResponse{
		Game: game.Games[params.GameID],
	})
}

// GetGameInfo godoc
// @Summary Get game information
// @Description Get's the information for a game defined by UUID
// @Tags Game
// @Accept json
// @Produce json
// @Param game_id path string true "game uuid"
// @Success 200 {object} GameInfoResponse
// @Failure 400 {object} ErrorResponse
// @Router /game/info/{game_id} [get]
func getGameInfo(c *gin.Context) {

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

	board := game.Games[params.GameID]

	i := &GameInfoResponse{
		Name:          board.Name,
		Status:        board.Status,
		SizeX:         board.SizeX,
		SizeY:         board.SizeY,
		CurrentPlayer: board.CurrentPlayer,
		Fow:           board.FOWDistance,
		Created:       board.Created,
		Started:       board.Started,
		GameOptions: GameOptionsResponse{
			FogOfWarPercent:       game.FogOfWarPercent,
			AttackRange:           game.AttackRange,
			PlayerRoundMoves:      game.PlayerRoundMoves,
			MaxRoundSeconds:       game.MaxRoundSeconds,
			MaxPlayerMoveDistance: game.MaxPlayerMoveDistance,
			PowerUpMax:            game.PowerUpMax,
		},
		GameEntities: GameEntitiesResponse{
			AliveCreep:   len(board.AliveCreep()),
			AlivePlayers: len(board.AlivePlayers()),
			PowerUps:     len(board.PowerUps),
		},
	}

	c.JSON(http.StatusOK, &i)
}

// GetEvents godoc
// @Summary Get game events
// @Description Get's the events for a game defined by UUID
// @Tags Game
// @Accept json
// @Produce json
// @Param game_id path string true "game uuid"
// @Success 200 {object} GameEventsResponse
// @Failure 400 {object} ErrorResponse
// @Router /game/events/{game_id} [get]
func getEvents(c *gin.Context) {

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

	c.JSON(http.StatusOK, &GameEventsResponse{
		Events: game.Games[params.GameID].Events,
	})
}

// GetScores godoc
// @Summary Get game scores
// @Description Get's the scores for a game defined by UUID
// @Tags Game
// @Accept json
// @Produce json
// @Param game_id path string true "game uuid"
// @Success 200 {object} GameScoresResponse
// @Failure 400 {object} ErrorResponse
// @Router /game/scores/{game_id} [get]
func getScores(c *gin.Context) {

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

	s := &GameScoresResponse{}

	for _, player := range game.Games[params.GameID].Players {
		s.Scores = append(s.Scores, &PlayerScore{
			Name:          player.Name,
			Score:         player.Score,
			DamageDealt:   player.DamageDealt,
			DamageTaken:   player.DamageTaken,
			CreepKilled:   player.CreepKilled,
			PlayersKilled: player.PlayersKilled,
		})
	}

	c.JSON(http.StatusOK, &s)
}

// StartGame godoc
// @Summary Start a game
// @Description Starts a game
// @Tags Game
// @Accept json
// @Produce json
// @Param game_id path string true "game uuid"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Router /game/start/{game_id} [put]
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

	game := game.Games[params.GameID]
	if len(game.Players) <= 0 {
		c.JSON(http.StatusBadRequest, &ErrorResponse{
			Message: `cannot start game without players`,
		})
		return
	}

	// start the game
	go game.Run()

	c.JSON(http.StatusOK, &StatusResponse{
		Success: true,
	})
}

// StopGame godoc
// @Summary Stop a game
// @Description Stop a game
// @Tags Game
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param game_id path string true "game uuid"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Router /game/stop/{game_id} [put]
func stopGame(c *gin.Context) {

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

	board := game.Games[params.GameID]
	board.Status = game.BoardStatusFinished

	c.JSON(http.StatusOK, &StatusResponse{
		Success: true,
	})
}

// JoinGame godoc
// @Summary Join a player to a game
// @Description Joins a player to an existing game
// @Tags Game
// @Accept json
// @Produce json
// @Param data body	JoinPlayerRequest true "Join Request"
// @Success 200 {object} StatusResponse
// @Failure 403 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Router /game/join [post]
func joinGame(c *gin.Context) {

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

	board := game.Games[params.GameID]
	if board.Status != game.BoardStatusNew {
		c.JSON(http.StatusForbidden, &ErrorResponse{
			Message: `game is not accepting new players`,
		})
		return
	}

	for _, p := range game.Games[params.GameID].Players {
		if params.PlayerID == p.ID {
			c.JSON(http.StatusForbidden, &ErrorResponse{
				Message: `player is already in the game`,
			})
			return
		}
	}

	var player storage.Player
	storage.Storage.Get().Where("UUID = ?", params.PlayerID).First(&player)

	gamePlayer := game.NewPlayer(&player)
	game.Games[params.GameID].JoinPlayer(gamePlayer)

	c.JSON(http.StatusOK, &StatusResponse{
		Success: true,
	})
}
