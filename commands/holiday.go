package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var holidays = []Holiday{
	{
		Start:    time.Date(2024, time.December, 16, 0, 0, 0, 0, time.UTC),
		End:      time.Date(2024, time.December, 16, 0, 0, 0, 0, time.UTC),
		IsOneDay: true,
		Message:  "Holiday: Victory Day",
	},
	{
		Start:    time.Date(2024, time.December, 25, 0, 0, 0, 0, time.UTC),
		End:      time.Date(2024, time.December, 25, 0, 0, 0, 0, time.UTC),
		IsOneDay: true,
		Message:  "Holiday: Christmas Day",
	},
	{
		Start:    time.Date(2025, time.February, 2, 0, 0, 0, 0, time.UTC),
		End:      time.Date(2025, time.February, 2, 0, 0, 0, 0, time.UTC),
		IsOneDay: true,
		Message:  "Holiday: Sharaswati Puja (Only classes will remain suspended)",
	},
	{
		Start:    time.Date(2025, time.February, 14, 0, 0, 0, 0, time.UTC),
		End:      time.Date(2025, time.February, 14, 0, 0, 0, 0, time.UTC),
		IsOneDay: true,
		Message:  "Holiday: *Shab-e-Barat",
	},
}

type Holiday struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	IsOneDay bool      `json:"isOneDay"`
	Message  string    `json:"message"`
}

var holidayHandler = Commnad{
	Trigger: "holiday",
	Command: &discordgo.ApplicationCommand{
		Name:        "holiday",
		Description: "Get the upcoming holiday details",
	},

	Handler: func(op *options) {
		// if true {
		// 	in_maintainance(op)
		// 	return
		// }

		now := time.Now()
		var nextHoliday Holiday

		for _, e := range holidays {
			if e.Start.After(now) {
				nextHoliday = e
				break
			}
		}

		var message string

		for i, e := range holidays {
			if e.IsOneDay {
				message += fmt.Sprintf("%d. %s - %s \n", i, e.Start.Format("02 Jan 2006"), e.Message)
			} else {
				totalDays := int((e.End.Sub(e.Start).Hours() / 24))
				message += fmt.Sprintf("%d. %s - %s ** %d days ** \n", i, e.Start.Format("02 Jan 2006"), e.Message, totalDays)
			}
		}

		message += "\nNext Holiday " + nextHoliday.Message + " start in " + formatDuration(nextHoliday.Start)
		message += "\n" + SUPPORT_STRING + "\n"
		op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
			},
		})
	},
}
