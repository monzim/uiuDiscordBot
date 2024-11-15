package main

import (
	"os"

	"github.com/joho/godotenv"
	db "github.com/monzim/uiuBot/database"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	err := godotenv.Load()
	if err != nil {
		log.Error().Err(err).Msg("Error loading .env file")
	}

	postgres, err := db.NewDatabaseConnection(
		// os.Getenv("LOG_DATABASE_URI"),
		os.Getenv("DATABASE_URI"),
	)
	if err != nil {
		log.Error().Err(err).Msg("Error initializing the database connection")

	}

	// AddExamsToDatabase()
	SendToUnizimTest(postgres)

}
