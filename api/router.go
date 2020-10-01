package api

import "github.com/gin-gonic/gin"

func SetupRouter() (r *gin.Engine) {
	r = gin.Default()

	r.GET("/ping", ping)

	gamer := r.Group("/games")
	{
		gamer.GET("/new", newGame)
		gamer.GET("/", allGames)
	}

	return
}
