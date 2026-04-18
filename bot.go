package discord

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session      *discordgo.Session
	Commands     map[string]Command
	DeveloperIDs []string

	dmMu       sync.Mutex
	dmChannels map[string]string // userID → channel ID cache
}

func NewBot(token string) *Bot {
	session, err := createSession(token)
	if err != nil {
		panic(err)
	}

	bot := &Bot{
		Session:    session,
		Commands:   make(map[string]Command),
		dmChannels: make(map[string]string),
	}

	bot.Session.AddHandler(bot.handleInteractionCreate)
	bot.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		slog.Info(fmt.Sprintf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator))
	})
	return bot
}

func (bot *Bot) Stop() {
	bot.Session.Close()
}

func createSession(Token string) (s *discordgo.Session, err error) {
	s, err = discordgo.New("Bot " + Token)
	if err != nil {
		slog.Error("Failed to create Discord session", "error", err)
		return
	}

	s.Identify.Intents = discordgo.IntentsGuildMessages

	err = s.Open()
	if err != nil {
		slog.Error("Failed to open websocket connection", "error", err)
		return
	}

	return
}

func (bot *Bot) AddCommands(cmds ...Command) {
	for _, cmd := range cmds {
		bot.Commands[cmd.Command.Name] = cmd
	}
}

func (bot *Bot) ClearCommands(guildID string) {
	cmds, _ := bot.Session.ApplicationCommands(bot.Session.State.User.ID, guildID)

	for _, cmd := range cmds {
		err := bot.Session.ApplicationCommandDelete(bot.Session.State.User.ID, guildID, cmd.ID)
		if err != nil {
			slog.Error("Failed to delete commands", "error", err)
		}
	}

	bot.Commands = make(map[string]Command)
}

func (bot *Bot) RegisterCommands(guildID string) {
	cmds := make([]*discordgo.ApplicationCommand, 0, len(bot.Commands))
	for _, cmd := range bot.Commands {
		cmds = append(cmds, &cmd.Command)
	}

	_, err := bot.Session.ApplicationCommandBulkOverwrite(
		bot.Session.State.User.ID,
		guildID,
		cmds,
	)
	if err != nil {
		slog.Error("Failed to register commands", "error", err)
	}
}

func (bot *Bot) handleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if cmd, ok := bot.Commands[i.ApplicationCommandData().Name]; ok {
		cmd.Handler(s, i)
	}
}

// channelIDByName returns the first text channel ID matching name across all cached guilds.
func (bot *Bot) channelIDByName(name string) (string, bool) {
	for _, guild := range bot.Session.State.Guilds {
		for _, ch := range guild.Channels {
			if ch.Type == discordgo.ChannelTypeGuildText && ch.Name == name {
				return ch.ID, true
			}
		}
	}
	return "", false
}

// dmChannelID returns the cached DM channel ID for a user, creating it if necessary.
func (bot *Bot) dmChannelID(userID string) (string, error) {
	bot.dmMu.Lock()
	defer bot.dmMu.Unlock()

	if id, ok := bot.dmChannels[userID]; ok {
		return id, nil
	}

	ch, err := bot.Session.UserChannelCreate(userID)
	if err != nil {
		return "", err
	}

	bot.dmChannels[userID] = ch.ID
	return ch.ID, nil
}

func (bot *Bot) SendChannelMessage(channelName string, message string) {
	chID, ok := bot.channelIDByName(channelName)
	if !ok {
		slog.Warn("Channel not found", "channel", channelName)
		return
	}
	if _, err := bot.Session.ChannelMessageSend(chID, message); err != nil {
		slog.Error("Failed to send channel message", "error", err, "channel", channelName)
	}
}

func (bot *Bot) SendEmbedMessage(channelName string, embed *discordgo.MessageEmbed) {
	chID, ok := bot.channelIDByName(channelName)
	if !ok {
		slog.Warn("Channel not found", "channel", channelName)
		return
	}
	if _, err := bot.Session.ChannelMessageSendEmbed(chID, embed); err != nil {
		slog.Error("Failed to send embed message", "error", err, "channel", channelName)
	}
}

func (bot *Bot) SendDeveloperMessage(message string) {
	for _, userID := range bot.DeveloperIDs {
		chID, err := bot.dmChannelID(userID)
		if err != nil {
			slog.Error("Failed to create DM channel", "error", err, "user", userID)
			continue
		}
		if _, err := bot.Session.ChannelMessageSend(chID, message); err != nil {
			slog.Error("Failed to send developer message", "error", err, "user", userID)
		}
	}
}

func (bot *Bot) SendDeveloperEmbed(embed *discordgo.MessageEmbed) {
	for _, userID := range bot.DeveloperIDs {
		chID, err := bot.dmChannelID(userID)
		if err != nil {
			slog.Error("Failed to create DM channel", "error", err, "user", userID)
			continue
		}
		if _, err := bot.Session.ChannelMessageSendEmbed(chID, embed); err != nil {
			slog.Error("Failed to send developer embed", "error", err, "user", userID)
		}
	}
}
