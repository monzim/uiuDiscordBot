package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

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

	logPg, err := db.NewDatabaseConnection(&db.DatabaseConfig{
		Host:     os.Getenv("LOG_DB_HOST"),
		Port:     os.Getenv("LOG_DB_PORT"),
		DBname:   os.Getenv("LOG_DB_NAME"),
		User:     os.Getenv("LOG_DB_USER"),
		Password: os.Getenv("LOG_DB_PASSWORD"),
		SSlMode:  os.Getenv("LOG_DB_SSL_MODE"),
	})

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

	go pingServerStatus(myBot)
	go statsPing(myBot)

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

func pingServerStatus(myBot *bot.Bot) {
	interval, err := time.ParseDuration(os.Getenv("STATUS_PING_INTERVAL_DURATION"))
	if err != nil {
		log.Error().Err(err).Msg("Error parsing the interval")
		interval = time.Minute * 60
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	startTime := time.Now()

	uptime := func() time.Duration {
		return time.Since(startTime).Round(time.Second)
	}

	for range ticker.C {
		myBot.Session.ChannelMessageSend(os.Getenv("STATUS_CHANNEL"), fmt.Sprintf("Hello, I'm still alive! Uptime: %v", uptime()))
	}
}

func statsPing(myBot *bot.Bot) {
	for {
		interval, err := time.ParseDuration(os.Getenv("USED_STATS_PING_INTERVAL_DURATION"))
		if err != nil {
			log.Error().Err(err).Msg("Error parsing the interval")
			interval = time.Minute * 140
		}

		startTime := time.Now().Add(-interval)
		myBot.LogServerStats()

		var userActivity models.UserActivity
		myBot.DB.Model(&models.UserActivity{}).Where("time_stamps > ?", startTime).Find(&userActivity)

		var message string

		if userActivity.CommandsExecuted > 0 {
			if userActivity.CommandsExecuted == 1 {
				message = fmt.Sprintf("### In the last %s, 1 command was executed", interval.String())
			} else if userActivity.CommandsExecuted > 10 {
				message = fmt.Sprintf("### Wow! In the last %s, %d commands were executed", interval.String(), userActivity.CommandsExecuted)
			} else {
				message = fmt.Sprintf("### In the last %s, %d commands were executed", interval.String(), userActivity.CommandsExecuted)
			}
		} else {
			message = fmt.Sprintf("### In the last %s, no commands were executed :(", interval.String())
		}

		myBot.Session.ChannelMessageSend(os.Getenv("STATUS_CHANNEL"), message)

		time.Sleep(interval)
	}
}

func handleUnexpectedPanics() {
	if r := recover(); r != nil {
		log.Error().Interface("panic", r).Msg("Unexpected panic")
		os.Exit(1)
	}
}
