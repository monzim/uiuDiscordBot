package commands

import (
	"github.com/bwmarrin/discordgo"
)

var handlerAuthor = Commnad{
	Trigger: "author",
	Command: &discordgo.ApplicationCommand{
		Name:        "author",
		Description: "Replies with the author's information",
	},

	Handler: func(op *options) {
		op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This bot was created by <@669529872644833290>\nhttps://monzim.com\n" + SUPPORT_STRING,
			},
		})
	},
}
