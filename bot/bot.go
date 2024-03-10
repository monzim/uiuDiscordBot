package bot

import (
	"encoding/json"

	"github.com/bwmarrin/discordgo"
	"github.com/monzim/uiuBot/commands"
	"github.com/monzim/uiuBot/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Bot struct {
	Session       *discordgo.Session
	RemoveCommand bool
	DB            *gorm.DB
	LogDb         *gorm.DB
}

func NewBot(token, guildID string, removeCommands bool, db *gorm.DB, logDb *gorm.DB) (*Bot, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		Session:       s,
		RemoveCommand: removeCommands,
		DB:            db,
		LogDb:         logDb,
	}, nil
}

func (b *Bot) Open() error {
	err := b.Session.Open()
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) Close() {
	b.Session.Close()
}

func (b *Bot) AddCommandHandlers() {
	b.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.GuildID == "" {
			log.Info().Msg("Ignoring DM interaction")

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "### Oops! This command is not supported in DM. Please use it in a server. Join our server: https://monzim.com/uiubot",
				},
			})

			go func() {
				b.LogDb.Create(&models.EventLog{
					EventType:        "DM",
					DM:               true,
					EventDescription: "User tried to interact with the bot in DM " + i.ApplicationCommandData().Name + " command",
				})

				data, err := json.Marshal(i)
				if err != nil {
					log.Error().Err(err).Msg("Error marshalling the message")
					return
				}

				user, err := json.Marshal(s.State.User)
				if err != nil {
					log.Error().Err(err).Msg("Error marshalling the message")
					return
				}

				b.LogDb.Create(&models.DMLog{
					UserID:   s.State.User.ID,
					Data:     data,
					UserData: user,
				})

			}()

			return
		}

		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			commands.HandleCommand(s, i, b.DB, b.LogDb)

		case discordgo.InteractionMessageComponent:
			commands.HandlerComponents(s, i, b.DB, b.LogDb)
		}
	})
}

func (b *Bot) ListCommands(guildID string) []*discordgo.ApplicationCommand {
	commands, err := b.Session.ApplicationCommands(b.Session.State.User.ID, guildID)
	if err != nil {
		log.Error().Err(err).Msg("Cannot get commands")
		return nil
	}

	return commands
}

func (b *Bot) RegisterCommands(commands []*discordgo.ApplicationCommand, guildID string) []*discordgo.ApplicationCommand {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	for i, v := range commands {
		cmd, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, guildID, v)
		log.Debug().Msgf("Registered command: %v", v.Name)

		if err != nil {
			log.Error().Err(err).Msgf("Cannot create '%v' command", v.Name)
		}

		registeredCommands[i] = cmd
	}

	return registeredCommands
}

func (b *Bot) RemoveCommands(commands []*discordgo.ApplicationCommand, guildID string) {
	log.Info().Msg("Removing commands...")

	for _, v := range commands {
		log.Warn().Msgf("Removing command: %v", v.Name)
		err := b.Session.ApplicationCommandDelete(b.Session.State.User.ID, guildID, v.ID)
		if err != nil {
			log.Error().Err(err).Msgf("Cannot delete '%v' command", v.Name)
		}
	}
}

func (b *Bot) LogServerStats() {
	guilds, err := b.Session.UserGuilds(100, "", "")
	if err != nil {
		log.Error().Err(err).Msg("Cannot get guilds")
		return
	}

	for _, guild := range guilds {
		g, err := b.Session.Guild(guild.ID)
		if err != nil {
			log.Error().Err(err).Msgf("Cannot get guild info for %v", guild.ID)
			return
		}

		b.LogDb.FirstOrCreate(&models.ServerStats{
			ServerID:      g.ID,
			MembersCount:  len(g.Members),
			ChannelsCount: len(g.Channels),
		})
	}
}
