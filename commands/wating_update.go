package commands

import "github.com/bwmarrin/discordgo"

func in_maintainance(op *options) {
	op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "This command is on maintainance mode\n" + SUPPORT_STRING,
		},
	})
}
