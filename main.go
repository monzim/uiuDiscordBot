package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/monzim/uiuBot/bot"
	"github.com/monzim/uiuBot/commands"
	db "github.com/monzim/uiuBot/database"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutting down or not")
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := godotenv.Load()
	if err != nil {
		log.Error().Err(err).Msg("Error loading .env file")
	}

	*BotToken = os.Getenv("TOKEN")
}

func main() {
	flag.Parse()

	postgres, err := db.NewDatabaseConnection(&db.DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBname:   os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		SSlMode:  os.Getenv("DB_SSL_MODE"),
	})
	if err != nil {
		log.Error().Err(err).Msg("Error initializing the database connection")

	}

	err = db.Migrate(postgres)
	if err != nil {
		log.Error().Err(err).Msg("Error migrating the database")
	}

	myBot, err := bot.NewBot(*BotToken, *GuildID, *RemoveCommands, postgres)
	if err != nil {
		log.Error().Err(err).Msg("Invalid bot parameters")
	}

	defer myBot.Close()

	// Open the bot session before registering commands
	err = myBot.Open()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot open the session")
	}

	myBot.AddCommandHandlers()

	log.Info().Msg("Adding commands...")
	registeredCommands := myBot.RegisterCommands(commands.GetCommands(), *GuildID)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Info().Msg("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Warn().Msg("Removing commands...")
		myBot.RemoveCommands(registeredCommands, *GuildID)
	}

	log.Info().Msg("Gracefully shutting down.")

}
