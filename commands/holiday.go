package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var holidays = []Holiday{

	{
		Start:    time.Date(2024, time.February, 21, 0, 0, 0, 0, time.UTC),
		End:      time.Date(2024, time.February, 21, 0, 0, 0, 0, time.UTC),
		IsOneDay: true,
		Message:  "Holiday: International Mother Language Day",
	},
	{
		Start:    time.Date(2024, time.February, 26, 0, 0, 0, 0, time.UTC),
		End:      time.Date(2024, time.February, 26, 0, 0, 0, 0, time.UTC),
		IsOneDay: true,
		Message:  "Holiday: *Shab-e-Barat",
	},
	{
		Start:    time.Date(2024, time.March, 17, 0, 0, 0, 0, time.UTC),
		End:      time.Date(2024, time.March, 17, 0, 0, 0, 0, time.UTC),
		IsOneDay: true,
		Message:  "Holiday: Birthday of Father of the Nation Bangabandhu Sheikh Mujibur Rahman",
	},
	{
		Start:    time.Date(2024, time.March, 26, 0, 0, 0, 0, time.UTC),
		End:      time.Date(2024, time.March, 26, 0, 0, 0, 0, time.UTC),
		IsOneDay: true,
		Message:  "Holiday: Independence Day",
	},
	{
		Start:    time.Date(2024, time.April, 4, 0, 0, 0, 0, time.UTC),
		End:      time.Date(2024, time.April, 14, 0, 0, 0, 0, time.UTC),
		IsOneDay: false,
		Message:  "Holiday:  Jumu’atul-Widaa/ *Shab-e-Qad’r /*Eid-ul-Fit’r/ Bangla New Year (for administration  only)",
	},
	{
		Start:    time.Date(2024, time.April, 4, 0, 0, 0, 0, time.UTC),
		End:      time.Date(2024, time.April, 15, 0, 0, 0, 0, time.UTC),
		IsOneDay: false,
		Message:  "Holiday:  Jumu’atul-Widaa/ *Shab-e-Qad’r /*Eid-ul-Fit’r/ Bangla New Year (for students only)",
	},
	{
		Start:    time.Date(2024, time.May, 1, 0, 0, 0, 0, time.UTC),
		End:      time.Date(2024, time.May, 1, 0, 0, 0, 0, time.UTC),
		IsOneDay: true,
		Message:  "Holiday: May Day",
	},
	{
		Start:    time.Date(2024, time.May, 22, 0, 0, 0, 0, time.UTC),
		End:      time.Date(2024, time.May, 22, 0, 0, 0, 0, time.UTC),
		IsOneDay: true,
		Message:  "Holiday: *Buddha Purnima",
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
