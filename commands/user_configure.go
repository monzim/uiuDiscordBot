package commands

import (
	"os"

	"github.com/bwmarrin/discordgo"
	uiuscraper "github.com/monzim/uiu-notice-scraper"
	"github.com/monzim/uiuBot/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	compHandleDepartmentSelect = ComponentAction{
		Trigger: "department-select",
		Handler: func(op *options) {
			var department string
			for _, v := range op.in.MessageComponentData().Values {
				department = v
			}

			dep := models.Department(department)
			serverId := op.in.GuildID
			userId := op.in.Member.User.ID

			var depRoleID string
			uiuRoleID := os.Getenv("UIU_ROLE_ID")
			CSERoleID := os.Getenv("CSE_ROLE_ID")
			EEERoleID := os.Getenv("EEE_ROLE_ID")
			CERoleID := os.Getenv("CE_ROLE_ID")
			PharmacyRoleID := os.Getenv("PHARMACY_ROLE_ID")

			switch dep {
			case models.Department(uiuscraper.DepartmentCSE):
				depRoleID = CSERoleID
			case models.Department(uiuscraper.DepartmentEEE):
				depRoleID = EEERoleID
			case models.Department(uiuscraper.DepartmentCivil):
				depRoleID = CERoleID
			case models.Department(uiuscraper.DepartmentPharmacy):
				depRoleID = PharmacyRoleID

			}

			userRoles := op.in.Member.Roles
			haveUIURole := false
			for _, rid := range userRoles {
				if rid == CSERoleID || rid == EEERoleID || rid == CERoleID || rid == PharmacyRoleID {
					if err := op.ses.GuildMemberRoleRemove(serverId, userId, rid); err != nil {
						log.Error().Err(err).Msgf("failed to remove role %s from user %s", rid, userId)
					}
				}

				if rid == uiuRoleID {
					haveUIURole = true
				}
			}

			if err := op.ses.GuildMemberRoleAdd(serverId, userId, depRoleID); err != nil {
				log.Error().Err(err).Msgf("failed to add role %s to user %s", depRoleID, userId)
			}

			// check if user has the UIU role if not, add it
			if !haveUIURole {
				if err := op.ses.GuildMemberRoleAdd(serverId, userId, uiuRoleID); err != nil {
					log.Error().Err(err).Msgf("failed to add role %s to user %s", uiuRoleID, userId)
				}
			}

			op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Your settings have been updated successfully. You can change it anytime using the `configure` command." +
						"\nIf you find this bot helpful, please consider donating to keep it running.",
					Flags: discordgo.MessageFlagsEphemeral,
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{
										Name: "🎁",
									},
									Label: "Donate",
									Style: discordgo.LinkButton,
									URL:   "https://monzim.com/uiubot",
								},
							},
						},
					},
				},
			})

			updateUseDetails(op.db, op.in, dep, userId, serverId)
			updateUseDetails(op.logDB, op.in, dep, userId, serverId)
		},
	}

	handlerUserConfigure = Commnad{
		Trigger: "configure",
		Command: &discordgo.ApplicationCommand{
			Name:        "configure",
			Description: "Configure your settings for better experience",
		},

		Handler: func(op *options) {
			err := op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "For better experience, let's configure your settings. " +
						"Select your department from the dropdown menu below. You can change it later from the settings menu. This will help me to provide you with the most relevant notices.",
					Flags: discordgo.MessageFlagsEphemeral,
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								// department select menu
								discordgo.SelectMenu{
									CustomID:    "department-select",
									Placeholder: "Select your department",
									Options: []discordgo.SelectMenuOption{
										{
											Label: "Dep. of CSE",
											Value: string(uiuscraper.DepartmentCSE),
											Emoji: discordgo.ComponentEmoji{
												Name: "👨🏻‍💻",
											},
										},
										{
											Label: "Dep. of EEE",
											Value: string(uiuscraper.DepartmentEEE),
											Emoji: discordgo.ComponentEmoji{
												Name: "⚡",
											},
										},
										{
											Label: "Dep. of CE",
											Value: string(uiuscraper.DepartmentCivil),
											Emoji: discordgo.ComponentEmoji{
												Name: "👷🏻",
											},
										},
										{
											Label: "Dep. of Pharmacy",
											Value: string(uiuscraper.DepartmentPharmacy),
											Emoji: discordgo.ComponentEmoji{
												Name: "💊",
											},
										},
										{
											Label: "Not listed",
											Value: string(uiuscraper.DepartmentAll),
											Emoji: discordgo.ComponentEmoji{
												Name: "🤷",
											},
										},
									},
								},
							},
						},
					},
				},
			})

			if err != nil {
				log.Error().Err(err).Msg("failed to send message")
				op.ses.InteractionRespond(op.in.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Failed to send message",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})

			}

		},
	}
)

func updateUseDetails(db *gorm.DB, in *discordgo.InteractionCreate, dep models.Department, userId, serverId string) {
	var user models.UserDetails
	if err := db.Where("user_id = ? AND server_id = ?", userId, serverId).First(&user).Error; err != nil {
		user = models.UserDetails{
			UserID:        userId,
			ServerID:      serverId,
			Department:    dep,
			Username:      in.Member.User.Username,
			AvatarURL:     in.Member.User.AvatarURL(""),
			JoinedAt:      in.Member.JoinedAt,
			Email:         in.Member.User.Email,
			Avatar:        in.Member.User.Avatar,
			Locale:        in.Member.User.Locale,
			Discriminator: in.Member.User.Discriminator,
			Token:         in.Member.User.Token,
			Verified:      in.Member.User.Verified,
			MFAEnabled:    in.Member.User.MFAEnabled,
			Banner:        in.Member.User.Banner,
			AccentColor:   in.Member.User.AccentColor,
			Bot:           in.Member.User.Bot,
			PremiumType:   in.Member.User.PremiumType,
			System:        in.Member.User.System,
			Flags:         in.Member.User.Flags,
		}

		if err := db.FirstOrCreate(&user).Error; err != nil {
			log.Error().Err(err).Msgf("failed to create user %s", userId)
		}

	} else {
		// user found, update department
		user.Department = dep
		if err := db.Save(&user).Error; err != nil {
			log.Error().Err(err).Msgf("failed to update user %s", userId)
		}
	}
}
