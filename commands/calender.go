package commands

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

var TRIMSTER_NAME = "Spring_2024"
var CALENDER_PATH = "public/Cal-Spring-2024.pdf"

var academyCalenderHandler = Commnad{
	Trigger: "academy-calender",
	Command: &discordgo.ApplicationCommand{
		Name:        "academy-calender",
		Description: "Replies with current academy calender",
	},

	Handler: func(op *options) {
		pf, err := os.Open(CALENDER_PATH)
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
				Content: "Here is the current academy calender pdf\n" + SUPPORT_STRING,
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
