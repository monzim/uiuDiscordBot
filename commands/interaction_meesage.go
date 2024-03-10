package commands

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type ComponentAction struct {
	Trigger string
	Handler func(op *options)
}

var (
	componentsHandlers = map[string]func(op *options){
		compHandleDepartmentSelect.Trigger: compHandleDepartmentSelect.Handler,
	}
)

func HandlerComponents(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB, logDb *gorm.DB) {
	if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
		h(&options{
			in:    i,
			ses:   s,
			db:    db,
			logDB: logDb,
		})
	}
}
