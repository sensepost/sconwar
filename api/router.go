package api

import "github.com/gin-gonic/gin"

// SetupRouter configures the HTTP routes we have
func SetupRouter() (r *gin.Engine) {
	r = gin.Default()

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
		}

		player := api.Group("/player")
		{
			player.GET("/")
			player.POST("/register")
			player.POST("/move")
			player.POST("/attack")
		}
	}

	return
}
