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

// @host localhost
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
			games.GET("/new", newGame)
			games.GET("/", allGames)
			games.GET("/get/:game_id", getGame)

			games.POST("/join", joinGame)
		}

		player := api.Group("/player")
		{
			player.GET("/")
			player.POST("/register")
		}
	}

	return
}