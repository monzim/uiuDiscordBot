package commands

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/monzim/uiuBot/models"
)

var handleSendMessage = Commnad{
	Trigger: "send-message",
	Command: &discordgo.ApplicationCommand{
		Name:        "send-message",
		Description: "This will send a message to you",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message",
				Description: "Input your message",
				Required:    true,
			},

			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "time",
				Description: "Choose time unit",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "minutes",
						Value: "minutes",
					},
					{
						Name:  "hours",
						Value: "hours",
					},
					{
						Name:  "days",
						Value: "days",
					},
				},
			},

			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "value",
				Description: "Input the value",
				Required:    true,
			},

			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "Select the channel",
				Required:    true,
			},
		},
	},

	Handler: func(op *options) {
		if !hasRole(op.in.Member.Roles, os.Getenv("ADMIN_ROLE_ID")) {
			op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You don't have permission to run this command",
				},
			})

			return
		}

		input := op.in.ApplicationCommandData().Options

		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(input))
		for _, opt := range input {
			optionMap[opt.Name] = opt
		}

		message := optionMap["message"].StringValue()
		timeUnit := optionMap["time"].StringValue()
		value := optionMap["value"].StringValue()
		channel := optionMap["channel"].ChannelValue(op.ses)

		sendTime := timeToDuration(timeUnit, value)
		res := op.db.Create(&models.CronMessage{
			Message:   message,
			ChannelID: channel.ID,
			UserID:    op.in.Member.User.ID,
			SendTime:  sendTime,
		})

		if res.Error != nil {
			op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to save schedule message. Please try again later.",
				},
			})
		}

		op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Message scheduled for %s in %s", sendTime.Format(time.RFC1123), channel.Mention()),
			},
		})

	},
}

func timeToDuration(timeUnit string, value string) time.Time {
	var duration time.Duration
	after, err := strconv.Atoi(value)
	if err != nil {
		duration = time.Minute * 10
	}

	switch timeUnit {
	case "minutes":
		duration = time.Minute * time.Duration(after).Abs()
	case "hours":
		duration = time.Hour * time.Duration(after).Abs()
	case "days":
		duration = time.Hour * 24 * time.Duration(after).Abs()
	}

	return time.Now().Add(duration)
}

func hasRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}

	return false
}
