package commands

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/monzim/uiuBot/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type options struct {
	in    *discordgo.InteractionCreate
	ses   *discordgo.Session
	db    *gorm.DB
	logDB *gorm.DB
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

func updateUserActivity(db *gorm.DB, userID string) {
	var userActivity models.UserActivity
	if err := db.FirstOrCreate(&userActivity, models.UserActivity{UserID: userID,
		ServerID: userID,
	}).Error; err != nil {
		return
	}

	db.Model(&userActivity).Updates(models.UserActivity{
		UserID:           userID,
		ServerID:         userID,
		CommandsExecuted: userActivity.CommandsExecuted + 1,
		LastActivity:     time.Now().String(),
	})
}

func HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB, logDb *gorm.DB) {

	go func() {
		updateUserActivity(logDb, i.Member.User.ID)
	}()

	go func() {
		logDb.Create(&models.EventLog{
			ServerID:         i.GuildID,
			EventType:        "command executed",
			EventDescription: i.ApplicationCommandData().Name,
		})

		// add or update user details
		var user models.UserDetails
		serverId := i.GuildID

		if err := logDb.FirstOrCreate(&user, models.UserDetails{
			ServerID:      serverId,
			UserID:        i.Member.User.ID,
			Username:      i.Member.User.Username,
			AvatarURL:     i.Member.User.AvatarURL(""),
			JoinedAt:      i.Member.JoinedAt,
			Email:         i.Member.User.Email,
			Avatar:        i.Member.User.Avatar,
			Locale:        i.Member.User.Locale,
			Discriminator: i.Member.User.Discriminator,
			Token:         i.Member.User.Token,
			Verified:      i.Member.User.Verified,
			MFAEnabled:    i.Member.User.MFAEnabled,
			Banner:        i.Member.User.Banner,
			AccentColor:   i.Member.User.AccentColor,
			Bot:           i.Member.User.Bot,
			PremiumType:   i.Member.User.PremiumType,
			System:        i.Member.User.System,
			Flags:         i.Member.User.Flags,
		}).Error; err != nil {
			return
		}

		logDb.Model(&user).Updates(models.UserDetails{
			ServerID:      serverId,
			UserID:        i.Member.User.ID,
			Username:      i.Member.User.Username,
			AvatarURL:     i.Member.User.AvatarURL(""),
			JoinedAt:      i.Member.JoinedAt,
			Email:         i.Member.User.Email,
			Avatar:        i.Member.User.Avatar,
			Locale:        i.Member.User.Locale,
			Discriminator: i.Member.User.Discriminator,
			Token:         i.Member.User.Token,
			Verified:      i.Member.User.Verified,
			MFAEnabled:    i.Member.User.MFAEnabled,
			Banner:        i.Member.User.Banner,
			AccentColor:   i.Member.User.AccentColor,
			Bot:           i.Member.User.Bot,
			PremiumType:   i.Member.User.PremiumType,
			System:        i.Member.User.System,
			Flags:         i.Member.User.Flags,
		})

	}()

	if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		startTime := time.Now()

		go func() {
			h(&options{in: i, ses: s, db: db, logDB: logDb})

			params := ""
			for _, v := range i.ApplicationCommandData().Options {
				params += v.Name + ": " + v.Value.(string) + " "
			}

			logDb.Create(&models.CommandLog{
				ServerID:     i.GuildID,
				UserID:       i.Member.User.ID,
				Command:      i.ApplicationCommandData().Name,
				Parameters:   params,
				ResponseTime: time.Since(startTime).String(),
			})
		}()
	}
}
