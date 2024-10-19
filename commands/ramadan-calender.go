package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/monzim/uiuBot/utils"
)

const calenderUrl = "https://cdn-monzim-com.azureedge.net/public-com/public/ce8dcb90-deff-11ee-9ef7-43de33b800f1"

var handlerRamadanCalender = Commnad{
	Trigger: "ramadan-calender",
	Command: &discordgo.ApplicationCommand{
		Name:        "ramadan-calender",
		Description: "Replies with 2024 Ramadan Calender",
	},

	Handler: func(op *options) {
		if true {
			in_maintainance(op)
			return
		}

		op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: op.in.Member.User.Mention() + " Ramadan Calender 2024. May Allah bless us all. Ameen. " +
					"\n" + utils.SUPPORT_MESSAGE,

				Embeds: []*discordgo.MessageEmbed{
					{
						Title: "Ramadan Calender 2024",
						URL:   calenderUrl,
						Image: &discordgo.MessageEmbedImage{
							URL: calenderUrl,
						},
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "UIU Discord Bot",
							IconURL: utils.BOT_LOGO,
						},
					},
				},
			},
		})
	},
}
