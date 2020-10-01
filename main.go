package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/sensepost/sconwar/api"
	"github.com/sensepost/sconwar/game"
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	rand.Seed(time.Now().UnixNano())

	game.Setup()
	api.SetupRouter().Run(":8080")
}
