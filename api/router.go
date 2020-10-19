package api

import (
	"log"
	"net/http"
	"os"

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

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name token

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
			game.POST("/new", newGame)
			game.POST("/join", joinGame)
			game.GET("/detail/:game_id", TokenAuthMiddleWare(), getGameDetail)
			game.GET("/events/:game_id", getEvents)
			game.PUT("/start/:game_id", startGame)

			game.GET("/info/:game_id", getGameInfo)
			game.GET("/scores/:game_id", getScores)
		}

		player := api.Group("/player")
		{
			player.POST("/", getPlayer)
			player.POST("/register", TokenAuthMiddleWare(), registerPlayer)

			player.POST("/status", playerStatus)
			player.POST("/surroundings", playerSurrounding)
			player.POST("/inventory")
		}

		action := api.Group("/action")
		{
			action.POST("/attack", attackAction)
			action.POST("/move", moveAction)
			action.POST("/pickup", pickupAction)
			action.POST("/use", useAction)
		}

		meta := api.Group("/meta")
		{
			meta.GET("/types", getTypes)
		}
	}

	return
}

// TokenAuthMiddleWare is a middleware that expects an API token
// ref: https://sosedoff.com/2014/12/21/gin-middleware.html
func TokenAuthMiddleWare() gin.HandlerFunc {
	requiredToken := os.Getenv("API_TOKEN")

	if requiredToken == "" {
		log.Fatal("Please set an API_TOKEN environment variable")
	}

	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API Token required"})
			return
		}

		if token != requiredToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Token"})
			return
		}

		c.Next()
	}
}
