package api

import (
	"github.com/gin-gonic/gin"
	_ "github.com/sensepost/sconwar/docs" // import auto generated docs
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title Sconwar API
// @version 1.0
// @description This is the sconwar API documentation.

// @host localhost:8080
// @BasePath /api

// SetupRouter configures the HTTP routes we have
func SetupRouter() (r *gin.Engine) {
	r = gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ping", ping)

	ui := r.Group("/ui")
	{
		ui.GET("/state")
	}

	api := r.Group("/api")
	{

		games := api.Group("/games")
		{
			games.GET("/", allGames)
			games.GET("/new", newGame)
			games.POST("/join", joinGame)
			games.GET("/get/:game_id", getGame)
			games.PUT("/start/:game_id", startGame)
		}

		player := api.Group("/player")
		{
			player.POST("/status")
			player.POST("/register", registerPlayer)
		}
	}

	return
}
