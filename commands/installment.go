package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Installment struct {
	Index   int       `json:"index"`
	Date    time.Time `json:"date"`
	Message string    `json:"message"`
	Color   int       `json:"color"`
}

var installments = []Installment{
	{
		Index:   0,
		Message: "1st installment: A fine of Tk. 1,000/- will be imposed if 40% Tuition Fee and Trimester Fee is not paid within this date. Transportation fee, if applicable, must be paid in full (no installment)",
		Date:    time.Date(2024, time.February, 13, 0, 0, 0, 0, time.UTC),
		Color:   0xA5DD9B,
	},

	{
		Index:   1,
		Message: "2nd installment: A fine of Tk. 1,000/- will be imposed if 70% Tuition Fee and Trimester Fee is not paid within this date.",
		Date:    time.Date(2024, time.March, 12, 0, 0, 0, 0, time.UTC),
		Color:   0xF9F07A,
	},

	{
		Index:   2,
		Message: "3rd installment: A fine of Tk. 1,000/- will be imposed if 100% Tuition Fee and Trimester Fee is not paid within this date.",
		Date:    time.Date(2024, time.April, 22, 0, 0, 0, 0, time.UTC),
		Color:   0xF7418F,
	},
}

var installmentHandler = Commnad{
	Trigger: "installment",
	Command: &discordgo.ApplicationCommand{
		Name:        "installment",
		Description: "Get the upcoming installment details",
	},

	Handler: func(op *options) {
		now := time.Now()
		var nextPayment Installment

		for _, e := range installments {
			if e.Date.After(now) {
				nextPayment = e
				break
			}
		}

		var embeds []*discordgo.MessageEmbed
		for _, e := range installments {
			embeds = append(embeds, &discordgo.MessageEmbed{
				Title:       fmt.Sprintf("Installment %d - **%s**", e.Index+1, e.Date.Format("02 Jan 2006")),
				Description: e.Message,
				Color:       e.Color,
			})

		}

		op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please note the following payment due dates and payment requirements. " +
					"Next Payment: **" + nextPayment.Date.Format("02 Jan 2006") +
					" time left " + formatDuration(nextPayment.Date) + " " +
					SUPPORT_STRING,
				Embeds: embeds,
			},
		})
	},
}
