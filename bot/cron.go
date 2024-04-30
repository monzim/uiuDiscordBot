package bot

import (
	"os"
	"time"

	"github.com/monzim/uiuBot/models"
	"github.com/rs/zerolog/log"
)

func (b *Bot) CronSchedule() {
	interval, err := time.ParseDuration(os.Getenv("CRON_INTERVAL_DURATION"))
	if err != nil {
		log.Error().Err(err).Msg("Error parsing the interval")
		interval = time.Minute * 5
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Error().Msgf("Recovered from panic: %v", r)
				}
			}()

			b.sendScheduledMessages(interval)
		}()
	}
}

func (b *Bot) sendScheduledMessages(interval time.Duration) {
	var messages []models.CronMessage
	err := b.DB.Where("done = ? AND send_time < ?", false, time.Now().Add(interval)).Find(&messages).Error

	log.Info().Msgf("CRON Found %d messages to send", len(messages))

	if err != nil {
		log.Error().Err(err).Msg("Error fetching scheduled messages")
	}

	for _, message := range messages {
		if message.SendTime.Before(time.Now()) {
			discordMessage, err := b.Session.ChannelMessageSend(message.ChannelID, message.Message)
			if err != nil {
				log.Error().Err(err).Msg("Error sending scheduled message")
			}

			message.Done = true
			err = b.DB.Save(&message).Error

			if err != nil {
				log.Error().Err(err).Msg("Error updating scheduled message")
			}

			log.Info().Msgf("Message sent to channel %s with message id %s", message.ChannelID, discordMessage.ID)
		}
	}
}
