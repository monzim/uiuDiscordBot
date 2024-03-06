package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
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
		ping.Trigger:                   ping.Handler,
		examTime.Trigger:               examTime.Handler,
		upcomingExam.Trigger:           upcomingExam.Handler,
		installmentHandler.Trigger:     installmentHandler.Handler,
		holidayHandler.Trigger:         holidayHandler.Handler,
		handlerAuthor.Trigger:          handlerAuthor.Handler,
		handlerVersion.Trigger:         handlerVersion.Handler,
		academyCalenderHandler.Trigger: academyCalenderHandler.Handler,
	}
)

func GetCommands(db *gorm.DB) []*discordgo.ApplicationCommand {
	var departments []string
	var departmentOptions []*discordgo.ApplicationCommandOptionChoice

	res := db.Table("exams").Distinct("department").Pluck("department", &departments)
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("Error fetching departments")
		departmentOptions = []*discordgo.ApplicationCommandOptionChoice{
			{Name: "BSCSE", Value: "BSCSE"},
			{Name: "BSDS", Value: "BSDS"},
			{Name: "BSEEE", Value: "BSEEE"},
			{Name: "BBA", Value: "BBA"},
			{Name: "BBA in AIS", Value: "BBA in AIS"},
			{Name: "BSECO", Value: "BSECO"},
			{Name: "BSCE", Value: "BSCE"},
			{Name: "BSSEDS", Value: "BSSEDS"},
			{Name: "MSCSE", Value: "MSCSE"},
		}
	} else {

		departmentOptions = make([]*discordgo.ApplicationCommandOptionChoice, len(departments))
		for i, d := range departments {
			departmentOptions[i] = &discordgo.ApplicationCommandOptionChoice{
				Name:  strings.ToUpper(d),
				Value: d,
			}
		}

	}

	return []*discordgo.ApplicationCommand{
		ping.Command,
		// examTime.Command,
		{
			Name:        "exam-time",
			Description: "Get your exam time",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "department",
					Description: "Which department are you from?",
					Required:    true,
					Choices:     departmentOptions,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "course_code",
					Description: "Input your course code",
					Required:    true,
					MinLength:   &MIN_COURSE_LEN,
					MaxLength:   10,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "section",
					Description: "Input your section",
					Required:    true,
					MinLength:   &MIN_SEC_LEN,
					MaxLength:   MAX_SEC_LEN,
				},
			},
		},
		upcomingExam.Command,
		installmentHandler.Command,
		holidayHandler.Command,
		handlerAuthor.Command,
		handlerVersion.Command,
		academyCalenderHandler.Command,
	}
}

func HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB) {
	if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		h(&options{in: i, ses: s, db: db})
	}
}
