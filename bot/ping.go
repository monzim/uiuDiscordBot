package bot

import (
	"fmt"
	"os"
	"time"

	"github.com/monzim/uiuBot/models"
	"github.com/rs/zerolog/log"
)

func (b *Bot) PingServerStatus() {
	idk := os.Getenv("STATUS_PING_INTERVAL_DURATION")
	if idk == "" {
		return
	}

	interval, err := time.ParseDuration(idk)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing the interval")
		interval = time.Hour * 2
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	startTime := time.Now()

	uptime := func() time.Duration {
		return time.Since(startTime).Round(time.Second)
	}

	for range ticker.C {
		b.Session.ChannelMessageSend(os.Getenv("STATUS_CHANNEL"), fmt.Sprintf("Hello, I'm still alive! Uptime: %v", uptime()))
	}
}

func (b *Bot) SendServerStatsToChannel() {
	idk := os.Getenv("USED_STATS_PING_INTERVAL_DURATION")
	if idk == "" {
		return
	}

	interval, err := time.ParseDuration(idk)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing the interval")
		interval = time.Hour * 3
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	count := 0

	for {
		count++

		b.LogServerStats()
		currentTime := time.Now()

		if count == 1 {
			continue
		}

		var userActivity models.UserActivity
		b.DB.Where("time_stamps > ?", currentTime.Add(-interval)).Find(&userActivity)

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

		b.Session.ChannelMessageSend(os.Getenv("STATUS_CHANNEL"), message)
		time.Sleep(interval)
	}
}
