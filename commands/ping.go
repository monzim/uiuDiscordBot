package commands

import (
	"github.com/bwmarrin/discordgo"
)

var ping = Commnad{
	Trigger: "ping",
	Command: &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Replies with pong and the latency of the bot",
	},

	Handler: func(op *options) {
		op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Pong! " + op.in.Member.User.Mention() + " Latency: " + op.ses.HeartbeatLatency().String(),
			},
		})
	},
}
