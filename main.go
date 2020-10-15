package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/sensepost/sconwar/api"
	"github.com/sensepost/sconwar/game"
	"github.com/sensepost/sconwar/storage"
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	rand.Seed(time.Now().UnixNano())

	storage.InitDb()

	game.Setup()
	api.SetupRouter().Run("localhost:8080")
}
