package commands

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

var TRIMSTER_NAME = "Fall_2024"
var calendar_PATH = "public/Cal-Fall-2024.pdf"

var academycalendarHandler = Commnad{
	Trigger: "academy-calendar",
	Command: &discordgo.ApplicationCommand{
		Name:        "academy-calendar",
		Description: "Replies with current academic calendar",
	},

	Handler: func(op *options) {
		pf, err := os.Open(calendar_PATH)
		if err != nil {
			op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error: " + err.Error(),
				},
			})
		}

		op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Here is the current academic calendar for " + TRIMSTER_NAME + ":\n" + SUPPORT_STRING,
				Files: []*discordgo.File{
					{
						ContentType: "application/pdf",
						Name:        TRIMSTER_NAME + "_Cal-UIU_BOT.pdf",
						Reader:      pf,
					},
				},
			},
		})
	},
}
