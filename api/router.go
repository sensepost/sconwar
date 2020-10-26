package api

import (
	"errors"
	"net/http"
	"os"

	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/rs/zerolog/log"
	docs "github.com/sensepost/sconwar/docs" // import auto generated docs
	"github.com/sensepost/sconwar/game"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	ginprometheus "github.com/zsais/go-gin-prometheus"
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

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	game.InitMetrics()

	bp := os.Getenv("BASE_HOST")
	if len(bp) != 0 {
		docs.SwaggerInfo.Host = os.Getenv("BASE_HOST")
	}

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
			meta.GET("/scores", Throttle(1), getTotalScores)
		}
	}

	return
}

// TokenAuthMiddleWare is a middleware that expects an API token
// ref: https://sosedoff.com/2014/12/21/gin-middleware.html
func TokenAuthMiddleWare() gin.HandlerFunc {
	requiredToken := os.Getenv("API_TOKEN")

	if requiredToken == "" {
		log.Fatal().Err(errors.New(`Please set an API_TOKEN environment variable`)).Msg(`failed to read api token`)
	}

	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API token required"})
			return
		}

		if token != requiredToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid API token"})
			return
		}

		c.Next()
	}
}

// Throttle rate limits an endpoint for i per second
func Throttle(i int) gin.HandlerFunc {
	lim := rate.NewLimiter(rate.Limit(i), i)

	return func(c *gin.Context) {
		if lim.Allow() {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusTooManyRequests,
			gin.H{"error": "rate limit exceeded"})
	}
}
