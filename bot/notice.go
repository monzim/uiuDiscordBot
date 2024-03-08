package bot

import (
	"os"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	uiuscraper "github.com/monzim/uiu-notice-scraper"
	"github.com/monzim/uiuBot/models"
	"github.com/monzim/uiuBot/utils"
	"github.com/rs/zerolog/log"
)

func (b *Bot) SendNotices() {
	var mutex sync.Mutex

	interval, err := time.ParseDuration(os.Getenv("NOTICE_CHECK_INTERVAL_DURATION"))
	if err != nil {
		log.Error().Err(err).Msg("Error parsing the interval")
		interval = time.Minute * 60
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		go b.ScrapNotices()

		var latestNotices []models.Notice
		res := b.DB.Where("notified = ?", false).Order("date asc").Find(&latestNotices)

		nLen := len(latestNotices)
		log.Info().Msgf("Found %d notices to send", nLen)

		if nLen == 0 {
			continue
		}

		if res.Error != nil {
			log.Warn().Err(res.Error).Msg("Error fetching the latest notices")
			continue
		}

		// send the notice to the channel
		channel := os.Getenv("NOTICE_CHANNEL")

		for _, notice := range latestNotices {
			mutex.Lock()
			log.Info().Msgf("Sending notice %v to the channel", notice.ID)

			if len(notice.Title) > 256 {
				notice.Title = notice.Title[:256]
			}

			embed := &discordgo.MessageEmbed{
				Title:       notice.Title,
				URL:         notice.Link,
				Description: utils.SUPPORT_MESSAGE,
				Image:       &discordgo.MessageEmbedImage{URL: notice.Image},
				Color:       utils.GenColorCode(notice.Title),
				Timestamp:   notice.Date.Format(time.RFC3339),
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "Notification from UIU Discord Bot",
					IconURL: utils.BOT_LOGO,
				},
			}

			b.Session.ChannelMessageSendEmbed(channel, embed)
			notice.Notified = true

			tx := b.DB.Begin()

			if err := tx.Save(&notice).Error; err != nil {
				tx.Rollback()
				log.Err(err).Msgf("Error updating the notice with notified status %v", notice.ID)
			} else {
				tx.Commit()
			}

			mutex.Unlock()
		}
	}
}

func (b *Bot) ScrapNotices() {
	var latestNotice models.Notice
	if err := b.DB.Order("date desc").First(&latestNotice).Error; err != nil {
		log.Warn().Err(err).Msg("No latest notice found in the database")
	}

	notices := uiuscraper.ScrapUIU(&latestNotice.ID)
	log.Info().Msgf("Scrapped %d notices", len(notices))

	for _, notice := range notices {
		var n models.Notice = models.Notice{
			ID:    notice.ID,
			Title: notice.Title,
			Image: notice.Image,
			Date:  notice.Date,
			Link:  notice.Link,
		}

		if err := b.DB.FirstOrCreate(&n).Error; err != nil {
			log.Error().Err(err).Msg("Error creating the notice")
		}
	}

	log.Info().Msg("Notices scrapped successfully")
}
