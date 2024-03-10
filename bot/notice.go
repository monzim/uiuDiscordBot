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
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Error().Msgf("Recovered from panic: %v", r)
				}
			}()

			b.ScrapNoticesByDepartment(
				uiuscraper.DepartmentAll,
				string(uiuscraper.AllowDomainUIU),
				string(uiuscraper.Notice_Site_UIU),
			)

			b.ScrapNoticesByDepartment(
				uiuscraper.DepartmentCSE,
				string(uiuscraper.AllowDomainCSE),
				string(uiuscraper.Notice_Site_CSE),
			)

			b.ScrapNoticesByDepartment(
				uiuscraper.DepartmentEEE,
				string(uiuscraper.AllowDomainEEE),
				string(uiuscraper.Notice_Site_EEE),
			)

		}()

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

		// send the notice to the channel_UIU
		channel_UIU := os.Getenv("NOTICE_CHANNEL")
		if channel_UIU == "" {
			log.Warn().Msg("No channel found to send the notice")
			continue
		}

		channel_CSE := os.Getenv("NOTICE_CHANNEL_CSE")
		if channel_CSE == "" {
			log.Warn().Msg("No channel found to send the notice")
			continue
		}

		channel_EEE := os.Getenv("NOTICE_CHANNEL_EEE")
		if channel_EEE == "" {
			log.Warn().Msg("No channel found to send the notice")
			continue
		}

		for _, notice := range latestNotices {
			mutex.Lock()
			log.Info().Msgf("Sending notice %v to the channel", notice.ID)

			if len(notice.Title) > 256 {
				notice.Title = notice.Title[:256]
			}

			uiuRoleID := os.Getenv("UIU_ROLE_ID")
			CSERoleID := os.Getenv("CSE_ROLE_ID")
			EEERoleID := os.Getenv("EEE_ROLE_ID")
			BBARoleID := os.Getenv("BBA_ROLE_ID")

			var mentionRoleID string
			switch notice.Department {
			case models.Department(uiuscraper.DepartmentCSE):
				mentionRoleID = CSERoleID
			case models.Department(uiuscraper.DepartmentEEE):
				mentionRoleID = EEERoleID
			case models.Department(uiuscraper.DepartmentBBA):
				mentionRoleID = BBARoleID
			case models.Department(uiuscraper.DepartmentAll):
				mentionRoleID = uiuRoleID
			}

			mentionRoleID = "<@&" + mentionRoleID + ">"

			embed := &discordgo.MessageEmbed{
				Title:       notice.Title,
				URL:         notice.Link,
				Description: mentionRoleID + " " + utils.SUPPORT_MESSAGE,
				Image:       &discordgo.MessageEmbedImage{URL: notice.Image},
				Color:       utils.GenColorCode(notice.Title),
				Timestamp:   notice.Date.Format(time.RFC3339),
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "Notification from UIU Discord Bot",
					IconURL: utils.BOT_LOGO,
				},
			}

			if notice.Department == models.Department(uiuscraper.DepartmentCSE) {
				b.Session.ChannelMessageSendEmbed(channel_CSE, embed)
			} else if notice.Department == models.Department(uiuscraper.DepartmentEEE) {
				b.Session.ChannelMessageSendEmbed(channel_EEE, embed)
			} else {
				b.Session.ChannelMessageSendEmbed(channel_UIU, embed)
			}

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

func (b *Bot) ScrapNoticesByDepartment(dep uiuscraper.Department, allowDomain string, noticeSite string) {
	var latestNotice models.Notice
	if err := b.DB.Order("date desc").Where("department = ?", dep).First(&latestNotice).Error; err != nil {
		log.Warn().Err(err).Msgf("No latest notice found in the database for department %s", dep)
	}

	config := uiuscraper.NoticeScrapConfig{
		LastNoticeId: &latestNotice.ID,
		Department:   dep,
		AllowDomain:  string(allowDomain),
		NOTICE_SITE:  string(noticeSite),
	}

	notices := uiuscraper.ScrapNotice(&config)
	log.Info().Msgf("Scrapped %d notices for department %s", len(notices), dep)
	for _, notice := range notices {
		var n models.Notice = models.Notice{
			ID:         notice.ID,
			Title:      notice.Title,
			Image:      notice.Image,
			Date:       notice.Date,
			Link:       notice.Link,
			Department: models.Department(notice.Department),
		}

		if err := b.DB.FirstOrCreate(&n).Error; err != nil {
			log.Error().Err(err).Msg("Error creating the notice")
		}
	}

	log.Info().Msgf("Department of %s notices scrapped successfully", dep)
}
