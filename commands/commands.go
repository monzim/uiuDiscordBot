package commands

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type options struct {
	in  *discordgo.InteractionCreate
	ses *discordgo.Session
	db  *gorm.DB
}

type Commnad struct {
	Trigger string
	Command *discordgo.ApplicationCommand
	Handler func(op *options)
}

var (
	commandHandlers = map[string]func(op *options){
		ping.Trigger:               ping.Handler,
		examTime.Trigger:           examTime.Handler,
		upcomingExam.Trigger:       upcomingExam.Handler,
		installmentHandler.Trigger: installmentHandler.Handler,
		holidayHandler.Trigger:     holidayHandler.Handler,
		handlerAuthor.Trigger:      handlerAuthor.Handler,
		handlerVersion.Trigger:     handlerVersion.Handler,
	}
)

func GetCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		ping.Command,
		examTime.Command,
		upcomingExam.Command,
		installmentHandler.Command,
		holidayHandler.Command,
		handlerAuthor.Command,
		handlerVersion.Command,
	}
}

func HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB) {
	if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		h(&options{in: i, ses: s, db: db})
	}
}
