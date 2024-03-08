package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/monzim/uiuBot/models"
	"github.com/monzim/uiuBot/utils"
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

		var notices []models.Notice
		res := op.db.Where("title ILIKE ?", "%"+searchTerm+"%").Order("date asc").Find(&notices)
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
			op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: op.in.Member.User.Mention() +
						"** Oho! We couldn't find any notice with the keyword **" + searchTerm + "**" + "\n** " + SUPPORT_STRING,
				},
			})

			return
		}

		var embeds []*discordgo.MessageEmbed
		for _, notice := range notices {
			notice.Title = strings.Title(notice.Title)

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

		op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: op.in.Member.User.Mention() +
					" " + "We found " +
					utils.Bold(fmt.Sprintf("%d", len(notices))) + " " + "matching " +
					" notices with the keyword " + utils.Bold(searchTerm) + " " +
					"in " + utils.Bold(elapsedTime.String()) + " seconds. " +
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
