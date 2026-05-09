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
	game.CreateGame(id, game.NewBoard(id, params.Name))

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
	for _, d := range game.ListGames() {
		g.Games = append(g.Games, &AllGamesGame{
			ID:     d.ID,
			Name:   d.Name,
			Status: int(d.StatusValue()),
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
// @Failure 401 {object} ErrorResponse
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

	board, _ := game.GetGame(params.GameID)
	c.JSON(http.StatusOK, &GameDetailResponse{
		Game: board,
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

	board, _ := game.GetGame(params.GameID)
	info := board.SnapshotGameInfo()

	i := &GameInfoResponse{
		Name:          info.Name,
		Status:        info.Status,
		SizeX:         info.SizeX,
		SizeY:         info.SizeY,
		CurrentPlayer: info.CurrentPlayer,
		Fow:           info.FOWDistance,
		Created:       info.Created,
		Started:       info.Started,
		GameOptions: GameOptionsResponse{
			FogOfWarPercent:       game.FogOfWarPercent,
			AttackRange:           game.AttackRange,
			PlayerRoundMoves:      game.PlayerRoundMoves,
			MaxRoundSeconds:       game.MaxRoundSeconds,
			MaxPlayerMoveDistance: game.MaxPlayerMoveDistance,
			PowerUpMax:            game.PowerUpMax,
		},
		GameEntities: GameEntitiesResponse{
			AliveCreep:   info.AliveCreep,
			AlivePlayers: info.AlivePlayers,
			PowerUps:     info.PowerUps,
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

	board, _ := game.GetGame(params.GameID)
	c.JSON(http.StatusOK, &GameEventsResponse{
		Events: board.SnapshotEvents(),
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

	board, _ := game.GetGame(params.GameID)
	s := &GameScoresResponse{}

	for _, player := range board.SnapshotPlayers() {
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

	board, _ := game.GetGame(params.GameID)
	if len(board.SnapshotPlayers()) <= 0 {
		c.JSON(http.StatusBadRequest, &ErrorResponse{
			Message: `cannot start game without players`,
		})
		return
	}
	if board.StatusValue() == game.BoardStatusRunning {
		c.JSON(http.StatusOK, &StatusResponse{Success: true})
		return
	}
	if board.StatusValue() == game.BoardStatusFinished {
		c.JSON(http.StatusBadRequest, &ErrorResponse{
			Message: `game has already finished`,
		})
		return
	}

	// start the game
	go board.Run()

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
// @Failure 401 {object} ErrorResponse
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

	board, _ := game.GetGame(params.GameID)
	if board.StatusValue() == game.BoardStatusFinished {
		c.JSON(http.StatusOK, &StatusResponse{Success: true})
		return
	}
	board.SetStatus(game.BoardStatusFinished)

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

	board, _ := game.GetGame(params.GameID)

	var player storage.Player
	storage.Storage.Get().Where("UUID = ?", params.PlayerID).First(&player)

	gamePlayer := game.NewPlayer(&player)
	if err := board.JoinPlayerIfOpen(gamePlayer); err != nil {
		c.JSON(http.StatusForbidden, &ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &StatusResponse{
		Success: true,
	})
}
