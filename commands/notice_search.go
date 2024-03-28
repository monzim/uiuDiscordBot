package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	uiuscraper "github.com/monzim/uiu-notice-scraper"
	"github.com/monzim/uiuBot/models"
	"github.com/monzim/uiuBot/utils"
	"gorm.io/gorm"
)

var (
	MIN_SEARCH_TERM_LEN = 3
	MAX_SEARCH_TERM_LEN = 30
)

var handlerNoticeSearch = Commnad{
	Trigger: "search-notice",
	Command: &discordgo.ApplicationCommand{
		Name:        "search-notice",
		Description: "Search notice with any keyword",
		Options: []*discordgo.ApplicationCommandOption{

			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "keyword",
				Description: "Search keyword",
				MinLength:   &MIN_SEARCH_TERM_LEN,
				MaxLength:   MAX_SEARCH_TERM_LEN,
				Required:    true,
			},

			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "department",
				Description: "Department",
				Required:    false,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "UIU", Value: uiuscraper.DepartmentAll},
					{Name: "BSCSE", Value: uiuscraper.DepartmentCSE},
					{Name: "BSEEE", Value: uiuscraper.DepartmentEEE},
				},
			},
		},
	},

	Handler: func(op *options) {
		startTime := time.Now()

		input := op.in.ApplicationCommandData().Options

		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(input))
		for _, opt := range input {
			optionMap[opt.Name] = opt
		}

		searchTerm := optionMap["keyword"].StringValue()
		dep := optionMap["department"]
		if dep == nil {
			dep = &discordgo.ApplicationCommandInteractionDataOption{
				Name:  "department",
				Value: "",
			}
		}

		department := dep.Value.(string)
		var notices []models.Notice

		var res *gorm.DB

		if department == "" || department == string(uiuscraper.DepartmentAll) {
			res = op.db.Where("title ILIKE ?", "%"+searchTerm+"%").
				Order("date asc").Find(&notices)
		} else {
			res = op.db.Where("title ILIKE ?", "%"+searchTerm+"%").
				Where("department = ?", department).
				Order("date asc").Find(&notices)
		}

		if res.Error != nil {
			op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: utils.ERROR_MESSAGE,
				},
			})

			return

		}

		if len(notices) == 0 {
			if department == "" || department == string(uiuscraper.DepartmentAll) {
				department = "UIU"
			}

			op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: op.in.Member.User.Mention() +
						"Oho! We couldn't find any notice with the keyword " + utils.Bold(searchTerm) +
						" with the department of " + utils.Bold(department) +
						". Please try again with a different keyword or department. " +
						". Query take " + utils.Bold(time.Since(startTime).String()) + ". " +
						SUPPORT_STRING,
				},
			})

			return
		}

		var embeds []*discordgo.MessageEmbed
		for _, notice := range notices {
			var title string
			if notice.Department != models.Department(uiuscraper.DepartmentAll) {
				title += "Dep. of " + string(notice.Department) + " - "
			}

			description := utils.ConstructDescription("", notice.Summary)

			title += notice.Title
			embed := &discordgo.MessageEmbed{
				Title:       strings.Title(title),
				URL:         notice.Link,
				Description: description,
				Image:       &discordgo.MessageEmbedImage{URL: notice.Image},
				Color:       utils.GenColorCode(notice.Title),
				Timestamp:   notice.Date.Format(time.RFC3339),
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "Notification from UIU Discord Bot",
					IconURL: utils.BOT_LOGO,
				},
			}

			embeds = append(embeds, embed)
		}

		elapsedTime := time.Since(startTime)

		// create a user channel to send the message
		channel, err := op.ses.UserChannelCreate(op.in.Member.User.ID)
		if err != nil {
			op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: utils.ERROR_MESSAGE,
				},
			})

			return
		}

		if department == "" || department == string(uiuscraper.DepartmentAll) {
			department = "UIU"
		}

		op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: op.in.Member.User.Mention() +
					" " + "We found " +
					utils.Bold(fmt.Sprintf("%d", len(notices))) + " " + "matching " +
					" notices with the keyword " + utils.Bold(searchTerm) +
					" with the department of " + utils.Bold(department) +
					" in " + utils.Bold(elapsedTime.String()) + " seconds. " +
					"Copied to your DM. " +
					"\n" + SUPPORT_STRING,
			},
		})

		for _, embed := range embeds {
			op.ses.ChannelMessageSendEmbed(channel.ID, embed)
		}

		op.ses.ChannelMessageSend(channel.ID, op.in.Member.User.Mention()+" "+utils.SUPPORT_MESSAGE)
	},
}
