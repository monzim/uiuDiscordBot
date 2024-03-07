package commands

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

var VERSION = "2.0.0-3"
var BUILD = time.Date(2024, time.March, 7, 01, 2, 17, 0, time.UTC)

var handlerVersion = Commnad{
	Trigger: "version",
	Command: &discordgo.ApplicationCommand{
		Name:        "version",
		Description: "Replies with the bot version",
	},

	Handler: func(op *options) {
		op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: op.in.Member.User.Mention() +
					" Version: " + VERSION + " Build time: " + BUILD.Format("02 Jan 2006 15:04:05") + " Latency: " + op.ses.HeartbeatLatency().String() + "\n" + SUPPORT_STRING,
			},
		})
	},
}
