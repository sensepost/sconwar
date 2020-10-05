package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
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
	r.Use(cors.Default())

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	ui := r.Group("/ui")
	{
		ui.GET("/state")
	}

	api := r.Group("/api")
	{

		game := api.Group("/game")
		{
			game.GET("/", allGames)
			game.GET("/new", newGame)
			game.POST("/join", joinGame)
			game.GET("/get/:game_id", getGame)
			game.PUT("/start/:game_id", startGame)
		}

		player := api.Group("/player")
		{
			player.POST("/register", registerPlayer)

			player.POST("/status", playerStatus)
			player.POST("/surroundings", playerSurrounding)
			player.POST("/inventory")
		}

		action := api.Group("/action")
		{
			action.POST("/attack", attackAction)
			action.POST("/move", moveAction)
			action.POST("/pickup")
		}
	}

	return
}
