package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/monzim/uiuBot/bot"
	"github.com/monzim/uiuBot/commands"
	db "github.com/monzim/uiuBot/database"
	"github.com/monzim/uiuBot/models"
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

	log.Info().Msg("Initializing the bot...")
	err := godotenv.Load()
	if err != nil {
		log.Error().Err(err).Msg("Error loading .env file")
	}

	*BotToken = os.Getenv("TOKEN")
}

func main() {
	defer handleUnexpectedPanics()

	log.Info().Msg("Starting the bot...")
	flag.Parse()

	postgres, err := db.NewDatabaseConnection(os.Getenv("DATABASE_URI"))
	if err != nil {
		log.Error().Err(err).Msg("Error initializing the database connection")

	}

	logPg, err := db.NewDatabaseConnection(os.Getenv("LOG_DATABASE_URI"))
	if err != nil {
		log.Error().Err(err).Msg("Error initializing the database connection 2")
	}

	err = db.Migrate(postgres)
	if err != nil {
		log.Error().Err(err).Msg("Error migrating the database")
	}

	err = db.Migrate(logPg)
	if err != nil {
		log.Error().Err(err).Msg("Error migrating the database 2")
	}

	myBot, err := bot.NewBot(*BotToken, *GuildID, *RemoveCommands, postgres, logPg)
	if err != nil {
		log.Error().Err(err).Msg("Invalid bot parameters")
	}

	defer myBot.Close()

	err = myBot.Open()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot open the session")
	}

	go myBot.LogServerStats()
	go func() {
		myBot.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
			if m.Author.ID == s.State.User.ID {
				return
			}

			jn, err := json.Marshal(m)
			if err != nil {
				log.Error().Err(err).Msg("Error marshalling the message")
				return
			}

			logPg.Create(models.MessageLog{
				ID:        m.ID,
				ServerID:  m.GuildID,
				UserID:    m.Author.ID,
				ChannelID: m.ChannelID,
				Message:   m.Content,
				Data:      jn,
			})
		})
	}()

	go myBot.PingServerStatus()
	go myBot.SendNotices()
	go myBot.CronSchedule()
	// go myBot.SendServerStatsToChannel()

	// list all commands
	myBot.ListCommands(*GuildID)

	myBot.AddCommandHandlers()

	log.Info().Msg("Adding commands...")
	registeredCommands := myBot.RegisterCommands(commands.GetCommands(myBot.DB), *GuildID)

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

func handleUnexpectedPanics() {
	if r := recover(); r != nil {
		log.Error().Interface("panic", r).Msg("Unexpected panic")
		os.Exit(1)
	}
}
