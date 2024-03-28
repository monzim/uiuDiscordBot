package commands

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/monzim/uiuBot/models"
	"github.com/monzim/uiuBot/utils"
)

var (
	MIN_COURSE_LEN = 3
	MIN_SEC_LEN    = 1
	MAX_SEC_LEN    = 3
	SUPPORT_STRING = utils.SUPPORT_MESSAGE
)

var examTime = Commnad{
	Trigger: "exam-time",
	Command: &discordgo.ApplicationCommand{
		Name:        "exam-time",
		Description: "Get your exam time",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "department",
				Description: "Which department are you from?",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "BSCSE", Value: "BSCSE"},
					{Name: "BSDS", Value: "BSDS"},
					{Name: "BSEEE", Value: "BSEEE"},
					{Name: "BBA", Value: "BBA"},
					{Name: "BBA in AIS", Value: "BBA in AIS"},
					{Name: "BSECO", Value: "BSECO"},
					{Name: "BSCE", Value: "BSCE"},
					{Name: "BSSEDS", Value: "BSSEDS"},
					{Name: "MSCSE", Value: "MSCSE"},
				},
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

	Handler: func(op *options) {
		startTime := time.Now()

		input := op.in.ApplicationCommandData().Options

		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(input))
		for _, opt := range input {
			optionMap[opt.Name] = opt
		}

		department := optionMap["department"].StringValue()
		courseCode := optionMap["course_code"].StringValue()
		section := optionMap["section"].StringValue()

		var exams []models.Exam

		res := op.db.Where(models.Exam{
			Department: strings.ToLower(department),
			Section:    strings.ToLower(section),
		}).Where("course_code LIKE ?", "%"+strings.ToLower(strings.TrimSpace(courseCode))+"%").Find(&exams)

		if res.Error != nil {
			op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error: " + res.Error.Error(),
				},
			})

		} else {
			courseCode = strings.ToUpper(courseCode)
			section = strings.ToUpper(section)
			department = strings.ToUpper(department)

			if len(exams) == 0 {
				op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: op.in.Member.User.Mention() +
							" **I couldn't find any exam time for Dep. ** " +
							utils.Bold(department) + " with course code " + utils.Bold(courseCode) + " and section " + utils.Bold(section) +
							". You may have entered course code incorrectly. Here is an example: `2213` or `cse 2213`. " +
							"\n" + SUPPORT_STRING,

						Embeds: []*discordgo.MessageEmbed{
							// TODO: remove this embed when the final exam schedule is available
							&discordgo.MessageEmbed{
								Title:       "This time is for Midterm",
								Description: "Midterm is over. Please check your final exam schedule. If you have any question, feel free to ask.",
								Color:       utils.GenColorCode("Midterm"),
							},
						}},
				})
				return
			}

			var embeds []*discordgo.MessageEmbed

			for _, exam := range exams {
				exam.CourseTitle = strings.Title(exam.CourseTitle)
				exam.CourseCode = strings.ToUpper(exam.CourseCode)
				exam.Section = strings.ToUpper(exam.Section)

				embeds = append(embeds, &discordgo.MessageEmbed{
					Title: exam.CourseTitle + " (" + strings.ToUpper(exam.Department) + ")",
					Color: utils.GenColorCode(exam.CourseCode),
					Description: "**" + exam.CourseCode + "**\n" +
						"Section **" + exam.Section + "**" + "     Faculty **" + exam.Teacher + "\n**\n**" + exam.ExamDate + " at " + exam.ExamTime + "**\n" + "Room " + exam.Room + "\n",
				})
			}

			// TODO: remove this embed when the final exam schedule is available
			embeds = append(embeds, &discordgo.MessageEmbed{
				Title:       "This time is for Midterm",
				Description: "Midterm is over. Please check your final exam schedule. If you have any question, feel free to ask.",
				Color:       utils.GenColorCode("Midterm"),
			})

			embeds[len(embeds)-1].Footer = &discordgo.MessageEmbedFooter{
				Text:    "Help Us Make a Difference",
				IconURL: "https://res.cloudinary.com/monzim/image/upload/v1688984685/download_kh1syl.png",
			}

			op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: op.in.Member.User.Mention() + " " + SUPPORT_STRING,
					Embeds:  embeds,
				},
			})
		}

		go func() {
			elapsed := time.Since(startTime)

			res := op.logDB.Create(&models.ExamTimeLog{
				ServerID:     op.in.GuildID,
				UserID:       op.in.Member.User.ID,
				Department:   department,
				CourseCode:   courseCode,
				Section:      section,
				ResponseTime: elapsed.String(),
			})

			if res.Error != nil {
				op.ses.ChannelMessageSend(op.in.ChannelID, "Error: "+res.Error.Error())
			}
		}()
	},
}
