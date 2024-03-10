package commands

import (
	"sort"

	"github.com/bwmarrin/discordgo"
	"github.com/monzim/uiuBot/utils"
)

type lol struct {
	Name        string
	Description string
}

var allCommands = []lol{
	{Name: ping.Trigger, Description: ping.Command.Description},
	{Name: "exam-time", Description: "Get your exam time information"},
	{Name: upcomingExam.Trigger, Description: upcomingExam.Command.Description},
	{Name: installmentHandler.Trigger, Description: installmentHandler.Command.Description},
	{Name: holidayHandler.Trigger, Description: holidayHandler.Command.Description},
	{Name: handlerAuthor.Trigger, Description: handlerAuthor.Command.Description},
	{Name: handlerVersion.Trigger, Description: handlerVersion.Command.Description},
	{Name: academyCalenderHandler.Trigger, Description: academyCalenderHandler.Command.Description},
	{Name: handlerNoticeSearch.Trigger, Description: handlerNoticeSearch.Command.Description},
	{Name: handlerUserConfigure.Trigger, Description: handlerUserConfigure.Command.Description},
}

var helpHandler = Commnad{
	Trigger: "help",
	Command: &discordgo.ApplicationCommand{
		Name:        "help",
		Description: "Replies with a list of available commands",
	},

	Handler: func(op *options) {
		sort.Slice(allCommands, func(i, j int) bool {
			return allCommands[i].Name < allCommands[j].Name
		})

		var message string
		for _, c := range allCommands {
			message += "`" + c.Name + "`: " + c.Description + "\n"
		}

		op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message + "\n" + utils.SUPPORT_MESSAGE,
			},
		})
	},
}
