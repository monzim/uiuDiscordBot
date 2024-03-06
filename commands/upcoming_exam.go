package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type NextExam struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Day       string    `json:"day"`
	Type      string    `json:"type"`
	TimeLeft  string    `json:"timeLeft"`
	Color     int       `json:"color"`
}

func formatDuration(inputTime time.Time) string {
	duration := time.Since(inputTime)
	absDuration := duration
	if duration < 0 {
		absDuration = -duration
	}

	days := int(absDuration.Hours() / 24)
	hours := int(absDuration.Hours()) % 24
	minutes := int(absDuration.Minutes()) % 60
	seconds := int(absDuration.Seconds()) % 60

	return fmt.Sprintf("%d days %d hours %d minutes %d seconds", days, hours, minutes, seconds)
}

func formatDateTime(inputTime time.Time) string {
	layout := "Monday, 2 Jan, 2006"
	formattedDateTime := inputTime.Format(layout)

	return formattedDateTime
}

var upcomingExams = []NextExam{
	{

		StartDate: time.Date(2024, time.March, 9, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, time.March, 16, 0, 0, 0, 0, time.UTC),
		Type:      "Mid-Term Exam",
		Color:     0x00ff00,
	},
	{

		StartDate: time.Date(2024, time.May, 12, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, time.May, 20, 0, 0, 0, 0, time.UTC),
		Type:      "Final Exam",
		Color:     0xff0000,
	},
}

var upcomingExam = Commnad{
	Trigger: "upcoming-exam",
	Command: &discordgo.ApplicationCommand{
		Name:        "upcoming-exam",
		Description: "Replies with the upcoming exam information",
	},

	Handler: func(op *options) {

		var embeds []*discordgo.MessageEmbed
		for _, e := range upcomingExams {
			embeds = append(embeds, &discordgo.MessageEmbed{
				Title:  e.Type + " from " + formatDateTime(e.StartDate) + " to " + formatDateTime(e.EndDate),
				Color:  e.Color,
				Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Time left: %s", formatDuration(e.StartDate))},
			})

		}

		op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Upcoming bamboo 🎋\n" + SUPPORT_STRING,
				Embeds:  embeds,
			},
		})
	},
}
